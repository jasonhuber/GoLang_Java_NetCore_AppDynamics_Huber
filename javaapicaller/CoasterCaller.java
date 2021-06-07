import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.concurrent.TimeUnit;


public class CoasterCaller {

	public static void main(String[] args) {
		while (true)
		{
			getcoasters();
			try {
				TimeUnit.SECONDS.sleep(20);
			} catch (InterruptedException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
			}
		}	}

	public static void getcoasters()
	{
//		https://mkyong.com/webservices/jax-rs/restfull-java-client-with-java-net-url/
				 try {

				        URL url = new URL("http://192.168.86.246:9090/coasters");
				        HttpURLConnection conn = (HttpURLConnection) url.openConnection();
				        conn.setRequestMethod("GET");
				        conn.setRequestProperty("Accept", "application/json");

				        if (conn.getResponseCode() != 200) {
				            throw new RuntimeException("Failed : HTTP error code : "
				                    + conn.getResponseCode());
				        }

				        BufferedReader br = new BufferedReader(new InputStreamReader(
				            (conn.getInputStream())));

				        String output;
				        System.out.println("Output from Server .... \n");
				        while ((output = br.readLine()) != null) {
				            System.out.println(output);
				        }

				        conn.disconnect();

				      } catch (MalformedURLException e) {

				        e.printStackTrace();

				      } catch (IOException e) {

				        e.printStackTrace();

				      }

				    }
		
}

