using DotNetEnv;
using GetStream;
using GetStream.Models;

namespace ServerCSharp.Services;

public class StreamService
{
    public StreamClient Client { get; private set; } = null!;

    public void Init()
    {
        Env.Load();

        var apiKey = Env.GetString("STREAM_API_KEY");
        var apiSecret = Env.GetString("STREAM_API_SECRET");

        if (string.IsNullOrEmpty(apiKey) || string.IsNullOrEmpty(apiSecret))
        {
            throw new InvalidOperationException("Stream API credentials missing. Set STREAM_API_KEY and STREAM_API_SECRET in .env");
        }

        Client = new StreamClient(apiKey, apiSecret);
    }

    public async Task CreateStreamUserAsync(string userId, string fullname, string profilePic)
    {
        var request = new UpdateUsersRequest
        {
            Users = new Dictionary<string, UserRequest>
            {
                [userId] = new UserRequest
                {
                    ID = userId,
                    Name = fullname,
                    Image = profilePic
                }
            }
        };

        await Client.UpdateUsersAsync(request);
    }

    public string GenerateStreamToken(string userId)
    {
        return Client.CreateUserToken(userId);
    }
}
