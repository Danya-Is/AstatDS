import org.json.JSONObject;
import java.io.*;
import java.net.Socket;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Iterator;
import java.util.LinkedList;

class ClientSock extends Thread {
    private final Socket socket;
    public final String ip;
    private final BufferedReader in;
    private final BufferedWriter out;

    private static final String getIps = "get ips request";
    private static final String getKV  = "get kvs request";
    private static final String getHashIps = "get ips hash request";
    private static final String getHashKV  = "get kvs hash request";

    public ClientSock(Socket socket, String ip) throws IOException {
        this.ip = ip;
        this.socket = socket;
        in = new BufferedReader(new InputStreamReader(socket.getInputStream()));
        out = new BufferedWriter(new OutputStreamWriter(socket.getOutputStream()));
        this.start();
    }

    public void checkKV() throws IOException {
        JSONObject request = new JSONObject().put("Type", getHashKV)
                .put("IP", "127.0.0.1:" + Server.State.techPort)
                .put("Key", "")
                .put("Value", "");
        out.write(request.toString() + "\n");
        out.flush();
        String line = in.readLine();
        if (!line.equals(CheckSum.md5(Server.State.KV.toString()))) {
            out.write(request.put("Type", getKV).toString() + "\n");
            out.flush();
            line = in.readLine();
            updateKV(new JSONObject(line));
        }
    }

    public void updateKV(JSONObject data) {
        JSONObject copy = new JSONObject(Server.State.KV.toString());
        Iterator<String> itr = data.keys();
        while (itr.hasNext()) {
            String ind = itr.next();
            if (copy.opt(ind) != null) {
                if (!data.getJSONObject(ind).getString("value").equals(copy.getJSONObject(ind).getString("value"))) {
                    JSONObject x = copy.getJSONObject(ind);
                    JSONObject y = data.getJSONObject(ind);
                    Date t1 = null;
                    Date t2 = null;
                    try {
                        t1 = new SimpleDateFormat(Server.State.format).parse(x.getString("time"));
                        t2 = new SimpleDateFormat(Server.State.format).parse(y.getString("time"));
                    } catch (ParseException e) {
                        System.out.println(e.getMessage());
                    }
                    if (t1 != null && t1.after(t2)) {
                        copy.put(ind, data.get(ind));
                    }
                }
            } else {
                copy.put(ind, data.get(ind));
            }
        }
        Server.State.KV = new JSONObject(copy.toString());
        Server.update();
    }

    public void checkIps() throws IOException, InterruptedException {
        JSONObject request = new JSONObject().put("Type", getHashIps)
                .put("IP", "127.0.0.1:" + Server.State.techPort)
                .put("Key", "")
                .put("Value", "");
        out.write(request.toString() + "\n");
        out.flush();
        String line = in.readLine();
        if (!line.equals(CheckSum.md5(Server.State.Ips.toString()))) {
            out.write(request.put("Type", getIps).toString() + "\n");
            out.flush();
            String line2 = in.readLine();
            updateIps(new JSONObject(line2));
        }
    }

    public void updateIps(JSONObject data) {
        JSONObject copy = new JSONObject(Server.State.Ips.toString());
        Iterator<String> itr = data.keys();
        while (itr.hasNext()) {
            String ind = itr.next();
            if (copy.opt(ind) != null) {
                if (!data.getJSONObject(ind).getString("status").equals(copy.getJSONObject(ind).getString("status"))) {
                    JSONObject x = copy.getJSONObject(ind);
                    JSONObject y = data.getJSONObject(ind);
                    Date t1 = null;
                    Date t2 = null;
                    try {
                        t1 = new SimpleDateFormat(Server.State.format).parse(x.getString("time"));
                        t2 = new SimpleDateFormat(Server.State.format).parse(y.getString("time"));
                    } catch (ParseException e) {
                        System.out.println(e.getMessage());
                    }
                    if (t1 != null && t1.after(t2)) {
                        copy.put(ind, data.get(ind));
                    }
                }
            } else {
                copy.put(ind, data.get(ind));
                WorkDispatcher.addNode(new JSONObject().put(ind, "{}"));
            }
        }
        Server.State.Ips = new JSONObject(copy.toString());
        Server.update();
    }

    @Override
    public void run() {
        try {
            while (true) {
                checkIps();
                checkKV();
                Thread.sleep(1000);
            }
        } catch (IOException | InterruptedException e) {
            System.out.println(e.getMessage());
        }
    }
}

public class WorkDispatcher extends Thread {

    public static LinkedList<ClientSock> clientList = new LinkedList<>();

    public static void addNode(JSONObject node) {
        Iterator<String> itr = node.keys();
        while (itr.hasNext()) {
            String nod = itr.next();
            if (Integer.parseInt(nod.split(":")[1]) == Server.State.techPort) {break;}
            boolean flag = false;
            for (ClientSock ip: clientList) {
                if (ip.ip.equals(nod)) {
                    flag = true;
                    break;
                }
            }
            if (flag) {break;}
            try {
                clientList.add(new ClientSock(new Socket(nod.split(":")[0], Integer.parseInt(nod.split(":")[1])), nod));
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
    }

    @Override
    public void run() {
        addNode(Server.State.Ips);
    }
}
