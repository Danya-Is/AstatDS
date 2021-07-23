import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import org.json.JSONObject;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Iterator;

public class HomeHandler implements HttpHandler {

    private static final String getValue = "get value request";
    private static final String getNodes = "get nodes request";
    private static final String putValue = "put value request";

    @Override
    public void handle(HttpExchange httpExchange) throws IOException {
        if ("GET".equals(httpExchange.getRequestMethod())) {
            handleGetRequest(httpExchange);
        } else if ("PUT".equals(httpExchange.getRequestMethod())) {
            handlePostRequest(httpExchange);
        }
    }

    public void handleGetRequest(HttpExchange httpExchange) throws IOException {
        String request = "incorrect request\n";
        BufferedReader isr = new BufferedReader(new InputStreamReader(httpExchange.getRequestBody()));
        String line = isr.readLine();
        JSONObject data = new JSONObject(line);
        if (getValue.equals(data.get("Type"))) {
            request = Server.getValue(data.get("Key").toString()) + "\n";
        } else if (getNodes.equals(data.get("Type"))) {
            request = "";
            Iterator itr = Server.State.Ips.keys();
            while (itr.hasNext()) {
                String elem = (String) itr.next();
                request += elem + "\n";
            }
        }

        httpExchange.sendResponseHeaders(200, request.getBytes().length);
        OutputStream output = httpExchange.getResponseBody();
        output.write(request.getBytes());
        output.flush();
        httpExchange.close();
    }

    public void handlePostRequest(HttpExchange httpExchange) throws IOException {
        String request = "incorrect request";
        BufferedReader isr = new BufferedReader(new InputStreamReader(httpExchange.getRequestBody()));
        String line = isr.readLine();
        JSONObject data = new JSONObject(line);
        if (putValue.equals(data.get("Type"))) {
            Server.putValue(data.get("Key").toString(), data.get("Value").toString(), new SimpleDateFormat(Server.State.format).format(new Date()));
            request = "OK\n";
        }

        httpExchange.sendResponseHeaders(200, request.getBytes().length);
        OutputStream output = httpExchange.getResponseBody();
        output.write(request.getBytes());
        output.flush();
        httpExchange.close();
    }
}

