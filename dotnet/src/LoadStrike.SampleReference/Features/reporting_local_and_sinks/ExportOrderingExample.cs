using LoadStrike;
using LoadStrike.Contracts.Metrics;
using LoadStrike.CrossPlatform;
using LoadStrike.CrossPlatform.AzureEventHubs;
using LoadStrike.CrossPlatform.Delegates;
using LoadStrike.CrossPlatform.Http;
using LoadStrike.CrossPlatform.Kafka;
using LoadStrike.CrossPlatform.Nats;
using LoadStrike.CrossPlatform.PushDiffusion;
using LoadStrike.CrossPlatform.RabbitMq;
using LoadStrike.CrossPlatform.RedisStreams;
using Microsoft.Extensions.Configuration;
using OpenQA.Selenium.Chrome;
using Serilog;
using Serilog.Events;

namespace LoadStrike.SampleReference.Features.ReportingLocalAndSinks;

/// <summary>
/// /// This feature area focuses on local reports and sink registration so the reporting feature is the main thing being explained. ///
/// </summary>
public static class ExportOrderingExample
{
    public static object Build()
    {
        // Export Ordering is the feature being explained. The sample returns that focused object alongside a minimal run result.
        var featureBundle = BuiltInSinks();
        var result = LoadStrikeRunner
            .RegisterScenarios(Scenario("export-ordering"))
            .WithRunnerKey(OrdersWorkflowPlaceholders.RunnerKey)
            .WithTestSuite("orders-reference")
            .WithTestName("export-ordering")
            .WithoutReports()
            .Run();

        return new { featureBundle, Result = result };
    }

        public static class OrdersWorkflowPlaceholders

        {

            public const string RunnerKey = "runner_dummy_orders_reference";

            public const string OrdersApiBaseUrl = "https://orders.example.test";

            public const string KafkaBootstrapServers = "localhost:9092";

            public const string NatsServerUrl = "nats://localhost:4222";

            public const string RedisConnectionString = "localhost:6379,abortConnect=false";

            public const string RabbitHost = "localhost";

            public const string AzureEventHubsConnectionString =

                "Endpoint=sb://orders.example.test/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=dummy";

            public const string ReportFolder = "./artifacts/reports";

            public const string OrderTopic = "orders.created";

            public const string ExampleOrderNumber = "ORD-10001";

            public const string ExampleTenant = "demo-tenant";



            public static string ScenarioName(string suffix) => $"orders.{suffix}";

        }

        public sealed class OrdersReportingSink : LoadStrikeReportingSink

        {

            public string SinkName => "orders-sample-sink";



            public Task Init(LoadStrikeBaseContext context, IConfiguration infraConfig) =>

                Task.CompletedTask;



            public Task Start(LoadStrikeSessionStartInfo sessionInfo) => Task.CompletedTask;



            public Task SaveRealtimeStats(LoadStrikeScenarioStats[] stats) => Task.CompletedTask;



            public Task SaveRealtimeMetrics(LoadStrikeMetricStats metrics) => Task.CompletedTask;



            public Task SaveRunResult(LoadStrikeRunResult result) => Task.CompletedTask;



            public Task Stop() => Task.CompletedTask;



            public void Dispose() { }

        }

        public sealed class OrdersRuntimePolicy : ILoadStrikeRuntimePolicy

        {

            public Task<bool> ShouldRunScenario(string scenarioName) =>

                Task.FromResult(!scenarioName.EndsWith(".skip", StringComparison.Ordinal));



            public Task BeforeScenario(string scenarioName) => Task.CompletedTask;



            public Task AfterScenario(string scenarioName, LoadStrikeScenarioRuntime stats) =>

                Task.CompletedTask;



            public Task BeforeStep(string scenarioName, string stepName) => Task.CompletedTask;



            public Task AfterStep(string scenarioName, string stepName, LoadStrikeReply reply) =>

                Task.CompletedTask;

        }

        public sealed class OrdersWorkerPlugin : ILoadStrikeWorkerPlugin

        {

            public string PluginName => "orders-sample-plugin";



            public Task Init(LoadStrikeBaseContext context, IConfiguration infraConfig) =>

                Task.CompletedTask;



            public Task Start(LoadStrikeSessionStartInfo sessionInfo) => Task.CompletedTask;



            public Task<LoadStrikePluginData> GetData(LoadStrikeRunResult result)

            {

                var table = LoadStrikePluginDataTable.Create("Captured Orders");

                table.Rows.Add(

                    new Dictionary<string, object>(StringComparer.Ordinal)

                    {

                        ["Order Number"] = OrdersWorkflowPlaceholders.ExampleOrderNumber,

                        ["Tenant"] = OrdersWorkflowPlaceholders.ExampleTenant,

                    }

                );



                var data = LoadStrikePluginData.Create(PluginName);

                data.Hints.Add("Reference-only plugin payload for the sample repository.");

                data.Tables.Add(table);

                return Task.FromResult(data);

            }



            public Task Stop() => Task.CompletedTask;



            public Task Dispose() => Task.CompletedTask;

        }

            public static Task<LoadStrikeReply> PublishReplyAsync(LoadStrikeScenarioContext context) =>

                Task.FromResult<LoadStrikeReply>(PublishReply(context));

            public static LoadStrikeReply PublishReply(LoadStrikeScenarioContext context) =>

                LoadStrikeStep

                    .Run(

                        "publish-order",

                        context,

                        () =>

                            Task.FromResult(

                                LoadStrikeResponse.Ok(

                                    new Dictionary<string, object>(StringComparer.Ordinal)

                                    {

                                        ["orderNumber"] = OrdersWorkflowPlaceholders.ExampleOrderNumber,

                                        ["tenant"] = OrdersWorkflowPlaceholders.ExampleTenant,

                                    },

                                    statusCode: "201",

                                    sizeBytes: 128,

                                    message: "created",

                                    customLatencyMs: 4.5

                                )

                            )

                    )

                    .GetAwaiter()

                    .GetResult()

                    .AsReply();

            public static LoadStrikeScenario Scenario(string suffix = "publish") =>

                LoadStrikeScenario

                    .Create(OrdersWorkflowPlaceholders.ScenarioName(suffix), PublishReplyAsync)

                    .WithInit(context => Task.CompletedTask)

                    .WithClean(context => Task.CompletedTask)

                    .WithWeight(2)

                    .WithoutWarmUp()

                    .WithMaxFailCount(3)

                    .WithRestartIterationOnFail(true);

            public static LoadStrikeRunner Runner() =>

                LoadStrikeRunner

                    .Create()

                    .AddScenario(Scenario())

                    .WithTestSuite("orders-reference")

                    .WithTestName("sample-reference")

                    .WithSessionId("orders-sample-session")

                    .WithReportFolder(OrdersWorkflowPlaceholders.ReportFolder)

                    .WithReportingInterval(TimeSpan.FromSeconds(2.5))

                    .WithoutReports()

                    .Configure(context =>

                        context

                            .WithRunnerKey(OrdersWorkflowPlaceholders.RunnerKey)

                            .WithReportFileName("orders-reference")

                            .WithReportFormats(LoadStrikeReportFormat.Txt, LoadStrikeReportFormat.Html)

                            .WithLoggerConfig(() => new LoggerConfiguration())

                            .WithMinimumLogLevel(LogEventLevel.Warning)

                            .WithRuntimePolicies(new OrdersRuntimePolicy())

                            .WithRuntimePolicyErrorMode(LoadStrikeRuntimePolicyErrorMode.Continue)

                            .WithWorkerPlugins(new OrdersWorkerPlugin())

                            .WithReportingSinks(new OrdersReportingSink())

                    );

            public static object BuiltInSinks() =>

                new

                {

                    InfluxDb = new InfluxDbReportingSink(

                        new InfluxDbReportingSinkOptions

                        {

                            BaseUrl = "http://localhost:8086",

                            Organization = "demo",

                            Bucket = "orders",

                            Token = "dummy-influx-token",

                        }

                    ),

                    Datadog = new DatadogReportingSink(

                        new DatadogReportingSinkOptions

                        {

                            BaseUrl = "https://http-intake.logs.datadoghq.com",

                            ApiKey = "dummy-datadog-api-key",

                            ApplicationKey = "dummy-datadog-app-key",

                        }

                    ),

                    GrafanaLoki = new GrafanaLokiReportingSink(

                        new GrafanaLokiReportingSinkOptions

                        {

                            BaseUrl = "http://localhost:3100",

                            TenantId = "demo-tenant",

                        }

                    ),

                    Splunk = new SplunkReportingSink(

                        new SplunkReportingSinkOptions

                        {

                            BaseUrl = "https://splunk.example.test",

                            Token = "dummy-splunk-token",

                        }

                    ),

                    OtelCollector = new OtelCollectorReportingSink(

                        new OtelCollectorReportingSinkOptions { BaseUrl = "http://localhost:4318" }

                    ),

                    TimescaleDb = new TimescaleDbReportingSink(

                        new TimescaleDbReportingSinkOptions

                        {

                            ConnectionString =

                                "Host=localhost;Port=5432;Database=orders;Username=orders;Password=dummy-password",

                            Schema = "public",

                            TableName = "loadstrike_reporting_events",

                            MetricsTableName = "loadstrike_reporting_metrics",

                        }

                    ),

                };
}
