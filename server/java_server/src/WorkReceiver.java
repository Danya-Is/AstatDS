import org.json.JSONObject;

import java.io.*;
import java.net.ServerSocket;
import java.net.Socket;
import java.net.SocketTimeoutException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.LinkedList;

class ServerSock extends Thread {
    private final Socket socket;
    private final BufferedReader in;
    private final BufferedWriter out;

    private static final String getIps = "get ips request";
    private static final String getKV  = "get kvs request";
    private static final String getHashIps = "get ips hash request";
    private static final String getHashKV  = "get kvs hash request";

    private static final String activated = "activated";
    private static final String deprecated  = "deprecated";

    private static String createAnswer(String line) {
        switch (line) {
            case getIps: return Server.State.Ips.toString() + "\n";
            case getKV: return Server.State.KV.toString() + "\n";
            case getHashKV: return CheckSum.md5(Server.State.KV.toString()) + "\n";
            case getHashIps: return CheckSum.md5(Server.State.Ips.toString()) + "\n";
        }
        return "invalid\n";
    }

    public ServerSock(Socket socket) throws IOException {
        this.socket = socket;
        in = new BufferedReader(new InputStreamReader(socket.getInputStream()));
        out = new BufferedWriter(new OutputStreamWriter(socket.getOutputStream()));
        this.start();
    }

    @Override
    public void run() {
        try {
            while (true) {
                try {
                    String line = in.readLine();
                    JSONObject request = new JSONObject(line);
                    JSONObject ip = Server.State.Ips.optJSONObject(request.getString("IP"));
                    if (ip == null) {
                        Server.State.Ips.put(request.getString("IP"), new JSONObject().put("status", activated).put("time", new SimpleDateFormat(Server.State.format).format(new Date())));
                        WorkDispatcher.addNode(new JSONObject().put(request.getString("IP"), "{}"));
                        Server.update();
                    } else if (ip.getString("status").equals(deprecated)) {
                        ip.put("status", activated);
                        ip.put("time", new SimpleDateFormat(Server.State.format).format(new Date()));
                        Server.update();
                        WorkDispatcher.addNode(new JSONObject().put(request.getString("IP"), "{}"));
                    }
                    out.write(createAnswer(request.getString("Type")));
                    out.flush();
                } catch (SocketTimeoutException | NullPointerException e) {
                    System.out.println(e.getMessage());
                    socket.close();
                    return;
                }
            }
        } catch (IOException e) {
            System.out.println(e.getMessage());
        }
    }
}

public class WorkReceiver implements Runnable {

    public static LinkedList<ServerSock> serverList = new LinkedList<>();

    @Override
    public void run() {
        try {
            ServerSocket serverSocket = new ServerSocket(Server.State.techPort);
            while (true) {
                Socket socket = serverSocket.accept();
                serverList.add(new ServerSock(socket));
            }
        } catch (IOException e) {
            System.out.println(e.getMessage());
        }
    }
}