import org.json.JSONObject;
import java.io.*;
import java.util.Base64;
import java.util.Properties;

public class Recorder {

    public static boolean createFile() {
        Properties p = System.getProperties();
        Server.State.statePath = p.getProperty("user.dir");
        File file = new File(Server.State.statePath + "/" + "Store-" + Server.State.techPort + ".txt");
        boolean isExist = file.exists();
        if (!isExist) {
            try {
                file.createNewFile();
            } catch (IOException e) {
                System.out.println(e.getMessage());
            }
        }
        return isExist;
    }

    public static void loadFromDisk() throws IOException {
        FileInputStream fin = new FileInputStream(Server.State.statePath + "/" + "Store-" + Server.State.techPort + ".txt");
        byte[] buffer = new byte[fin.available()];
        fin.read(buffer, 0, buffer.length);
        fin.close();

        String text = new String(buffer);
        byte[] decodedBytes = Base64.getDecoder().decode(text.replace("\n", ""));
        String decodedString = new String(decodedBytes);
        JSONObject state = new JSONObject(decodedString);
        Server.State.hash = state.getString("Hash");
        Server.State.KV = state.getJSONObject("KV");
        Server.State.Ips = state.getJSONObject("Ips");
    }

    public static void writeToDisk() throws IOException {
        JSONObject state = new JSONObject();
        state
                .put("Hash", Server.State.hash)
                .put("KV", Server.State.KV)
                .put("Ips", Server.State.Ips);
        String encodedString = Base64.getEncoder().encodeToString(state.toString().getBytes());

        FileOutputStream fout = new FileOutputStream(Server.State.statePath + "/" + "Store-" + Server.State.techPort + ".txt", false);
        byte[] buffer = encodedString.getBytes();
        fout.write(buffer, 0, buffer.length);
        fout.close();
    }
}
