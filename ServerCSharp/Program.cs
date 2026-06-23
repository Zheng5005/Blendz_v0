using DotNetEnv;
using ServerCSharp.Endpoints;
using ServerCSharp.Middleware;
using ServerCSharp.Services;

Env.Load();

var builder = WebApplication.CreateBuilder(args);

var port = Env.GetString("PORT");
if (!string.IsNullOrEmpty(port))
{
    builder.WebHost.UseUrls($"http://0.0.0.0:{port}");
}

// Register services
builder.Services.AddSingleton<DatabaseService>();
builder.Services.AddSingleton<StreamService>();
builder.Services.AddSingleton<JwtService>();

var app = builder.Build();

// Initialize MongoDB and Stream connections at startup
var databaseService = app.Services.GetRequiredService<DatabaseService>();
databaseService.Init();

var streamService = app.Services.GetRequiredService<StreamService>();
streamService.Init();

// CORS — match Go server exactly
app.Use(async (context, next) =>
{
    context.Response.Headers["Access-Control-Allow-Origin"] = "https://blendz-v0.vercel.app";
    context.Response.Headers["Access-Control-Allow-Headers"] = "Content-Type, Authorization";
    context.Response.Headers["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE, OPTIONS";
    context.Response.Headers["Access-Control-Allow-Credentials"] = "true";

    if (context.Request.Method == "OPTIONS")
    {
        context.Response.StatusCode = StatusCodes.Status204NoContent;
        return;
    }

    await next();
});

// Auth middleware — protects routes by validating the JWT cookie
app.UseMiddleware<AuthMiddleware>();

// Map endpoints
app.MapAuthEndpoints();
app.MapUserEndpoints();
app.MapChatEndpoints();

app.Run();
