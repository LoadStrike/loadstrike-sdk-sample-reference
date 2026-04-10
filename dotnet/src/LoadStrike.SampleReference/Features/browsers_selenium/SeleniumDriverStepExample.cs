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

namespace LoadStrike.SampleReference.Features.BrowsersSelenium;

/// <summary>
/// /// This feature area keeps the scenario small so the Selenium driver step is the main thing being explained. ///
/// </summary>
public static class SeleniumDriverStepExample
{
    public static object Build()
    {
        // Selenium Driver Step is the feature being explained. The sample returns that focused object alongside a minimal run result.
        var browserStep = SeleniumSample();
        dynamic browser = browserStep;
        var result = LoadStrikeRunner
            .RegisterScenarios(browser.BrowserStep)
            .WithRunnerKey(OrdersWorkflowPlaceholders.RunnerKey)
            .WithTestSuite("orders-reference")
            .WithTestName("selenium-driver-step")
            .WithoutReports()
            .Run();

        return new { browserStep, Result = result };
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

            public static object SeleniumSample()

            {

                var options = new ChromeOptions();

                options.AddArgument("--headless=new");

                options.AddArgument("--disable-gpu");



                return new

                {

                    ChromeOptions = options,

                    BrowserStep = LoadStrikeScenario.Create(

                        OrdersWorkflowPlaceholders.ScenarioName("selenium"),

                        async context =>

                            (

                                await LoadStrikeStep.Run(

                                    "selenium-order-checkout",

                                    context,

                                    () =>

                                        Task.FromResult(

                                            LoadStrikeResponse.Ok("selenium", "200", 256, "selenium", 8.5)

                                        )

                                )

                            ).AsReply()

                    ),

                };

            }
}
