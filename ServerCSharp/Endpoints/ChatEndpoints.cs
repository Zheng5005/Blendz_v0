using ServerCSharp.Services;

namespace ServerCSharp.Endpoints;

public static class ChatEndpoints
{
    public static void MapChatEndpoints(this IEndpointRouteBuilder app)
    {
        app.MapGet("/api/chat/token", GetStreamToken);
    }

    private static IResult GetStreamToken(HttpContext context, StreamService streamService)
    {
        if (context.Items["userId"] is not string userId)
        {
            return Results.Unauthorized();
        }

        try
        {
            var token = streamService.GenerateStreamToken(userId);
            return Results.Ok(token);
        }
        catch
        {
            return Results.StatusCode(StatusCodes.Status500InternalServerError);
        }
    }
}
