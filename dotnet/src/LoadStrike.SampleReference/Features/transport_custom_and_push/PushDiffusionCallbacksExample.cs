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

namespace LoadStrike.SampleReference.Features.TransportCustomAndPush;

/// <summary>
/// /// This feature area focuses on delegate and push callback contracts for custom transport integrations. ///
/// </summary>
public static class PushDiffusionCallbacksExample
{
    public static object Build()
    {
        // Push Diffusion Callbacks is the feature being explained. The sample returns that focused object alongside a minimal run result.
        var endpoints = PushEndpoints();
        var result = LoadStrikeRunner
            .RegisterScenarios(Scenario("push-diffusion-callbacks"))
            .WithRunnerKey(OrdersWorkflowPlaceholders.RunnerKey)
            .WithTestSuite("orders-reference")
            .WithTestName("push-diffusion-callbacks")
            .WithoutReports()
            .Run();

        return new { endpoints, Result = result };
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

            public static TrackingPayloadBuilder TrackingBuilder()

            {

                var builder = new TrackingPayloadBuilder

                {

                    ContentType = "application/json",

                    MessagePayloadType = typeof(Dictionary<string, object>),

                };

                builder.Headers["X-Correlation-Id"] = OrdersWorkflowPlaceholders.ExampleOrderNumber;

                builder.Headers["X-Tenant"] = OrdersWorkflowPlaceholders.ExampleTenant;

                builder.SetBody(

                    $$"""{"trackingId":"{{OrdersWorkflowPlaceholders.ExampleOrderNumber}}","tenantId":"{{OrdersWorkflowPlaceholders.ExampleTenant}}"}"""

                );

                return builder;

            }

            public static object DelegateEndpoints() =>

                new

                {

                    Producer = new DelegateStreamEndpointDefinition

                    {

                        Name = "orders-delegate-producer",

                        Mode = TrafficEndpointMode.Produce,

                        TrackingField = TrackingFieldSelector.Parse("header:x-id"),

                        ProduceAsync = (request, cancellationToken) =>

                            Task.FromResult(

                                new ProducedMessageResult

                                {

                                    IsSuccess = true,

                                    TimestampUtc = DateTimeOffset.UtcNow,

                                    ResponsePayload = TrackingBuilder().Build(),

                                }

                            ),

                    },

                    Consumer = new DelegateStreamEndpointDefinition

                    {

                        Name = "orders-delegate-consumer",

                        Mode = TrafficEndpointMode.Consume,

                        TrackingField = TrackingFieldSelector.Parse("json:$.trackingId"),

                        ConsumeAsync = async (onMessage, cancellationToken) =>

                            await onMessage(

                                new ConsumedMessage

                                {

                                    TimestampUtc = DateTimeOffset.UtcNow,

                                    Payload = TrackingBuilder().Build(),

                                }

                            ),

                    },

                };

            public static object PushEndpoints() =>

                new

                {

                    Producer = new PushDiffusionEndpointDefinition

                    {

                        Name = "orders-push-producer",

                        Mode = TrafficEndpointMode.Produce,

                        TrackingField = TrackingFieldSelector.Parse("header:x-id"),

                        ServerUrl = "ws://push.loadstrike.test",

                        TopicPath = "/orders",

                        PublishAsync = (request, cancellationToken) =>

                            Task.FromResult(

                                new ProducedMessageResult

                                {

                                    IsSuccess = true,

                                    TimestampUtc = DateTimeOffset.UtcNow,

                                    ResponsePayload = TrackingBuilder().Build(),

                                }

                            ),

                    },

                    Consumer = new PushDiffusionEndpointDefinition

                    {

                        Name = "orders-push-consumer",

                        Mode = TrafficEndpointMode.Consume,

                        TrackingField = TrackingFieldSelector.Parse("json:$.trackingId"),

                        ServerUrl = "ws://push.loadstrike.test",

                        TopicPath = "/orders",

                        SubscribeAsync = async (onMessage, cancellationToken) =>

                            await onMessage(

                                new ConsumedMessage

                                {

                                    TimestampUtc = DateTimeOffset.UtcNow,

                                    Payload = TrackingBuilder().Build(),

                                }

                            ),

                    },

                };
}
