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

namespace LoadStrike.SampleReference.Features.ConfigAndCli;

/// <summary>
/// /// This feature area focuses on loading config, overriding config with run arguments, and documenting the expected config layout. ///
/// </summary>
public static class DummyConfigLayoutExample
{
    public static object Build()
    {
        // Dummy Config Layout keeps the file-based setup visible while the runtime path stays minimal.
        var configLayout = new
        {
            ConfigFile = "loadstrike.config.json",
            InfraFile = "loadstrike.infra.json",
            RunnerKey = OrdersWorkflowPlaceholders.RunnerKey,
        };

        var context = LoadStrikeRunner
            .RegisterScenarios(Scenario("config-layout"))
            .LoadConfig($"features/config_and_cli/{configLayout.ConfigFile}")
            .LoadInfraConfig($"features/config_and_cli/{configLayout.InfraFile}")
            .WithRunnerKey(configLayout.RunnerKey);

        return LoadStrikeRunner.Run(context);
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

            public static LoadStrikeContext ConfiguredContext() =>

                LoadStrikeRunner

                    .RegisterScenarios(Scenario())

                    .LoadConfig("features/config_and_cli/loadstrike.config.json")

                    .LoadInfraConfig("features/config_and_cli/loadstrike.infra.json")

                    .WithRunnerKey(OrdersWorkflowPlaceholders.RunnerKey)

                    .WithSessionId("orders-config-sample");
}
