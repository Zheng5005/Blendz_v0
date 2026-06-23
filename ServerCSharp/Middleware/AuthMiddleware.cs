using ServerCSharp.Services;
using ServerCSharp.Utils;

namespace ServerCSharp.Middleware;

public class AuthMiddleware
{
    private readonly RequestDelegate _next;

    public AuthMiddleware(RequestDelegate next)
    {
        _next = next;
    }

    public async Task InvokeAsync(HttpContext context, JwtService jwtService)
    {
        if (context.Request.Path.StartsWithSegments("/api/auth/signup") ||
            context.Request.Path.StartsWithSegments("/api/auth/login") ||
            context.Request.Path.StartsWithSegments("/api/auth/logout"))
        {
            await _next(context);
            return;
        }

        var token = context.Request.Cookies[CookieHelper.CookieName];

        if (string.IsNullOrEmpty(token))
        {
            context.Response.StatusCode = StatusCodes.Status401Unauthorized;
            await context.Response.WriteAsync("Missing or invalid token");
            return;
        }

        var userId = jwtService.ValidateToken(token);

        if (string.IsNullOrEmpty(userId))
        {
            context.Response.StatusCode = StatusCodes.Status401Unauthorized;
            await context.Response.WriteAsync("Invalid token or expired");
            return;
        }

        context.Items["userId"] = userId;
        await _next(context);
    }
}
