using System;
    using System.Collections.Generic;
    using System.Linq;
    using System.Threading.Tasks;
    using LoadStrike;
    using LoadStrike.Contracts.Metrics;
    using LoadStrike.CrossPlatform;
    using LoadStrike.CrossPlatform.Http;
    using Microsoft.Extensions.Configuration;
    using Serilog;
    using Serilog.Events;

    namespace LoadStrike.SampleReference.Methods.OwnerLoadStrikeRunner;

    /// <summary>Configure the NATS endpoint used by clustered execution.</summary>
    public static class WithNatsServerUrlMethodReference
    {
        private const string RunnerKey = "runner_dummy_orders_reference";

private static Task<LoadStrikeReply<Dictionary<string, object>>> CreateTypedOrderReplyAsync()
{
    return Task.FromResult(
        LoadStrikeResponse.Ok(
            new Dictionary<string, object>(StringComparer.Ordinal)
            {
                ["orderId"] = "ORD-10001",
            },
            statusCode: "200",
            sizeBytes: 128,
            message: "ok",
            customLatencyMs: 3.2
        )
    );
}

private static LoadStrikeScenario baselineScenario(string name = "orders.get-by-id")
{
    return LoadStrikeScenario
        .Create(name, ExecuteOrderGetAsync)
        .WithLoadSimulations(LoadStrikeSimulation.IterationsForConstant(1, 1))
        .WithoutWarmUp();
}

private static async Task<LoadStrikeReply> ExecuteOrderGetAsync(LoadStrikeScenarioContext context)
{
    var reply = await LoadStrikeStep.Run("get-order", context, CreateTypedOrderReplyAsync);
    return reply.AsReply();
}

private static LoadStrikeRunner baseRunner()
{
    return LoadStrikeRunner
        .Create()
        .AddScenario(baselineScenario())
        .Configure(context => context.WithRunnerKey(RunnerKey))
        .WithTestSuite("orders-reference")
        .WithTestName("orders-get-by-id")
        .WithoutReports();
}

private static LoadStrikeContext baseContext() => baseRunner().BuildContext();

private static LoadStrikeContext trackedContext()
{
    return LoadStrikeRunner
        .Create()
        .AddScenario(
            CrossPlatformScenarioConfigurator.Configure(
                baselineScenario("orders.tracked"),
                trackingConfiguration()
            )
        )
        .Configure(context => context.WithRunnerKey(RunnerKey))
        .WithoutReports()
        .BuildContext();
}

private static HttpEndpointDefinition httpSource() =>
    new()
    {
        Name = "orders-http-source",
        Mode = TrafficEndpointMode.Produce,
        TrackingField = TrackingFieldSelector.Parse("header:X-Correlation-Id"),
        Url = "https://orders.example.test/api/orders",
        Method = "GET",
        TrackingPayloadSource = HttpTrackingPayloadSource.Request,
        ResponseSource = HttpResponseSource.ResponseBody,
    };

private static HttpEndpointDefinition httpDestination() =>
    new()
    {
        Name = "orders-http-destination",
        Mode = TrafficEndpointMode.Consume,
        TrackingField = TrackingFieldSelector.Parse("json:$.trackingId"),
        GatherByField = TrackingFieldSelector.Parse("json:$.tenantId"),
        Url = "https://orders.example.test/api/order-events",
        Method = "GET",
        ResponseSource = HttpResponseSource.ResponseBody,
        ConsumeJsonArrayResponse = true,
        ConsumeArrayPath = "$.items",
    };

private static CrossPlatformTrackingConfiguration trackingConfiguration()
{
    return new CrossPlatformTrackingConfiguration
    {
        Source = httpSource(),
        Destination = httpDestination(),
        RunMode = TrackingRunMode.GenerateAndCorrelate,
        CorrelationTimeout = TimeSpan.FromSeconds(30),
        TimeoutSweepInterval = TimeSpan.FromSeconds(1),
        TimeoutBatchSize = 200,
        TimeoutCountsAsFailure = true,
        MetricPrefix = "orders_tracking",
        ExecuteOriginalScenarioRun = false,
        CorrelationStore = CorrelationStoreConfiguration.InMemory(),
    };
}

private static TrackingPayload buildTrackingPayload()
{
    var builder = new TrackingPayloadBuilder();
    builder.SetBody("{\"trackingId\":\"ord-1\"}");
    return builder.Build();
}

private static LoadStrikeRunResult runResult() => default!;

private static LoadStrikeScenarioStats scenarioStats() => default!;

private static LoggerConfiguration buildLogger() => new();

private static TempConfigPaths writeTempConfigFiles()
{
    return new TempConfigPaths("method-reference.loadstrike.config.json", "method-reference.loadstrike.infra.json");
}

private readonly record struct TempConfigPaths(string ConfigPath, string InfraPath);

private sealed class OrdersReportingSink : LoadStrikeReportingSink
{
    public string SinkName => "orders-sample-sink";
    public Task Init(LoadStrikeBaseContext context, IConfiguration infraConfig) => Task.CompletedTask;
    public Task Start(LoadStrikeSessionStartInfo sessionInfo) => Task.CompletedTask;
    public Task SaveRealtimeStats(LoadStrikeScenarioStats[] stats) => Task.CompletedTask;
    public Task SaveRealtimeMetrics(LoadStrikeMetricStats metrics) => Task.CompletedTask;
    public Task SaveRunResult(LoadStrikeRunResult result) => Task.CompletedTask;
    public Task Stop() => Task.CompletedTask;
    public void Dispose() { }
}

private sealed class OrdersRuntimePolicy : ILoadStrikeRuntimePolicy
{
    public Task<bool> ShouldRunScenario(string scenarioName) => Task.FromResult(true);
    public Task BeforeScenario(string scenarioName) => Task.CompletedTask;
    public Task AfterScenario(string scenarioName, LoadStrikeScenarioRuntime stats) => Task.CompletedTask;
    public Task BeforeStep(string scenarioName, string stepName) => Task.CompletedTask;
    public Task AfterStep(string scenarioName, string stepName, LoadStrikeReply reply) => Task.CompletedTask;
}

private sealed class OrdersWorkerPlugin : ILoadStrikeWorkerPlugin
{
    public string PluginName => "orders-sample-plugin";
    public Task Init(LoadStrikeBaseContext context, IConfiguration infraConfig) => Task.CompletedTask;
    public Task Start(LoadStrikeSessionStartInfo sessionInfo) => Task.CompletedTask;
    public Task<LoadStrikePluginData> GetData(LoadStrikeRunResult result) => Task.FromResult(LoadStrikePluginData.Create(PluginName));
    public Task Stop() => Task.CompletedTask;
    public Task Dispose() => Task.CompletedTask;
}

        /// <summary>Apply the primary string value shown in the sample method reference.</summary>
public static object? ApplyPrimaryValueExample()
{
    return LoadStrikeRunner.WithNatsServerUrl(baseContext(), "nats://localhost:4222");
}

/// <summary>Show the same method with a second concrete value.</summary>
public static object? ApplyAlternateValueExample()
{
    return LoadStrikeRunner.WithNatsServerUrl(baseContext(), "nats://demo-cluster:4222");
}
    }
