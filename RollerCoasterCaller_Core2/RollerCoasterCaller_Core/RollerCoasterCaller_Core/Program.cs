using System.Collections.Generic;
using System.Diagnostics;
using System.Net.Http;
using System.Threading.Tasks;
using Newtonsoft.Json;


namespace RollerCoasterCaller_Core
{
    public class Coaster
    {
        public string Id { get; set; }
        public string Name { get; set; }
        public string Manufacturer { get; set; }
        public string InPark { get; set; }
        public int Height { get; set; }
    }

    class Program
    {
        static async Task Main(string[] args)
        {
            while (true)
            {
                _ = Task.Delay(10000);
                await ExternalCaller();
            }   
        }
        static async Task ExternalCaller()
        {
            List<Coaster> reservationList = new List<Coaster>();
            using (var httpClient = new HttpClient())
            {
                using (var response = await httpClient.GetAsync("http://192.168.86.246:9090/coasters"))
                {
                    string apiResponse = await response.Content.ReadAsStringAsync();
                    reservationList = JsonConvert.DeserializeObject<List<Coaster>>(apiResponse);
                    Debug.WriteLine(apiResponse);
                    System.Console.Out.WriteLine(apiResponse);
                }
            }
        }
    }
}


 
