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

namespace LoadStrike.SampleReference.Features.ThresholdsAndMetrics;

/// <summary>
/// /// This feature area focuses on where thresholds attach, how custom metrics are registered, and how threshold results surface in the run output. ///
/// </summary>
public static class StepThresholdsExample
{
    public static object Build()
    {
        // Step Thresholds keeps the scenario setup focused on one behavior. Change the option below first when comparing suboptions.
        var scenario = ScenarioWithThreshold(
                                "step-threshold",
                                LoadStrikeThreshold.CreateStep(
                                    "publish-order",
                                    stats => stats.Ok.Request.Count >= 1,
                                    1,
                                    TimeSpan.FromSeconds(2)
                                )
                            );
        return LoadStrikeRunner
            .RegisterScenarios(scenario)
            .WithRunnerKey(OrdersWorkflowPlaceholders.RunnerKey)
            .WithTestSuite("orders-reference")
            .WithTestName("step-thresholds")
            .WithoutReports()
            .Run();
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

            public static LoadStrikeScenario ScenarioWithThreshold(

                string suffix,

                LoadStrikeThreshold threshold

            ) => Scenario(suffix).WithThresholds(threshold);

            public static LoadStrikeScenario ScenarioWithMetrics()

            {

                var counter = Metric.CreateCounter("orders_published_total", "count");

                var gauge = Metric.CreateGauge("orders_payload_bytes", "bytes");



                return LoadStrikeScenario

                    .Create(

                        OrdersWorkflowPlaceholders.ScenarioName("metrics"),

                        async context =>

                        {

                            counter.Add(1);

                            gauge.Set(128);

                            return (

                                await LoadStrikeStep.Run(

                                    "publish-order",

                                    context,

                                    () =>

                                        Task.FromResult(

                                            LoadStrikeResponse.Ok(

                                                "metric-captured",

                                                "200",

                                                128,

                                                "metric-captured",

                                                3.2

                                            )

                                        )

                                )

                            ).AsReply();

                        }

                    )

                    .WithInit(context =>

                    {

                        context.RegisterMetric(counter);

                        context.RegisterMetric(gauge);

                        return Task.CompletedTask;

                    });

            }
}
