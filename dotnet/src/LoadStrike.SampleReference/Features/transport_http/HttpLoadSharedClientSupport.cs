using System.Net;
using System.Net.Http.Json;
using LoadStrike;

namespace LoadStrike.SampleReference.Features.TransportHttp;

internal static class HttpLoadSharedClientSupport
{
    private static readonly HttpClient SharedHttpClient = CreateHttpClient();

    private static HttpClient CreateHttpClient()
    {
        var handler = new SocketsHttpHandler
        {
            PooledConnectionLifetime = TimeSpan.FromMinutes(5),
            PooledConnectionIdleTimeout = TimeSpan.FromMinutes(2),
            MaxConnectionsPerServer = 1000,
            AutomaticDecompression = DecompressionMethods.GZip | DecompressionMethods.Deflate
        };

        return new HttpClient(handler)
        {
            Timeout = TimeSpan.FromSeconds(15)
        };
    }

    public static async Task<string> PostOrderStatusCodeAsync(
        string orderId,
        string tenantId,
        CancellationToken cancellationToken = default)
    {
        var payload = new
        {
            orderId,
            tenantId,
            amount = 49.95m
        };

        using var response = await SharedHttpClient.PostAsJsonAsync(
            "https://api.example.com/orders",
            payload,
            cancellationToken);

        return ((int)response.StatusCode).ToString();
    }

    public static Task DisposeSharedHttpClientAsync()
    {
        SharedHttpClient.Dispose();
        return Task.CompletedTask;
    }
}
