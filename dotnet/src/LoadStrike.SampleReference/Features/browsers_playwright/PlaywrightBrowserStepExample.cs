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

namespace LoadStrike.SampleReference.Features.BrowsersPlaywright;

/// <summary>
/// /// This feature area keeps the scenario small so the Playwright browser step is the main thing being explained. ///
/// </summary>
public static class PlaywrightBrowserStepExample
{
    public static object Build()
    {
        // Playwright Browser Step is the feature being explained. The sample returns that focused object alongside a minimal run result.
        var browserStep = PlaywrightSample();
        dynamic browser = browserStep;
        var result = LoadStrikeRunner
            .RegisterScenarios(browser.BrowserStep)
            .WithRunnerKey(OrdersWorkflowPlaceholders.RunnerKey)
            .WithTestSuite("orders-reference")
            .WithTestName("playwright-browser-step")
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

            public static object PlaywrightSample() =>

                new

                {

                    BrowserStep = LoadStrikeScenario.Create(

                        OrdersWorkflowPlaceholders.ScenarioName("playwright"),

                        async context =>

                        {

                            using var playwright = await Microsoft.Playwright.Playwright.CreateAsync();

                            var browser = await playwright.Chromium.LaunchAsync(

                                new Microsoft.Playwright.BrowserTypeLaunchOptions { Headless = true }

                            );

                            await browser.CloseAsync();

                            return (

                                await LoadStrikeStep.Run(

                                    "playwright-order-checkout",

                                    context,

                                    () =>

                                        Task.FromResult(

                                            LoadStrikeResponse.Ok(

                                                "playwright",

                                                "200",

                                                512,

                                                "playwright",

                                                10

                                            )

                                        )

                                )

                            ).AsReply();

                        }

                    ),

                };
}
