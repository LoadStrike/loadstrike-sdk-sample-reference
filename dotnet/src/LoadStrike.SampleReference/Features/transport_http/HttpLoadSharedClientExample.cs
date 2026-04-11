using LoadStrike;

namespace LoadStrike.SampleReference.Features.TransportHttp;

/// <summary>
/// Reuse one shared HttpClient for HTTP load scenarios so connection pooling stays stable under sustained request volume.
/// </summary>
public static class HttpLoadSharedClientExample
{
    public static object Build()
    {
        try
        {
            var result = LoadStrikeRunner
                .RegisterScenarios(Scenario())
                .WithRunnerKey("runner_dummy_orders_reference")
                .WithTestSuite("orders-reference")
                .WithTestName("http-load-shared-client")
                .WithoutReports()
                .Run();

            return new
            {
                Note = "Dispose SharedHttpClient when the run ends.",
                Result = result
            };
        }
        finally
        {
            HttpLoadSharedClientSupport.DisposeSharedHttpClientAsync().GetAwaiter().GetResult();
        }
    }

    public static LoadStrikeScenario Scenario() =>
        LoadStrikeScenario.Create("http-load-shared-client", async context =>
        {
            context.ScenarioInstanceData["tenantId"] ??= "tenant-a";
            var orderId = $"ord-{context.InvocationNumber}";

            var step = await LoadStrikeStep.Run<string>("POST /orders", context, async () =>
            {
                var statusCode = await HttpLoadSharedClientSupport.PostOrderStatusCodeAsync(
                    orderId,
                    (string)context.ScenarioInstanceData["tenantId"]);

                return statusCode.StartsWith('2')
                    ? LoadStrikeResponse.Ok("shared-http-client", statusCode: statusCode)
                    : LoadStrikeResponse.Fail<string>(statusCode: statusCode);
            });

            return step.AsReply();
        })
        .WithLoadSimulations(LoadStrikeSimulation.Inject(25, TimeSpan.FromSeconds(1), TimeSpan.FromSeconds(30)));
}
