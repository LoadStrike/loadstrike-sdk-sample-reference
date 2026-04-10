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

namespace LoadStrike.SampleReference.Features.TransportOtherBrokers;

/// <summary>
/// /// This feature area focuses on the minimal endpoint options for NATS, RabbitMQ, Redis Streams, and Azure Event Hubs. ///
/// </summary>
public static class AzureEventHubsExample
{
    public static object Build()
    {
        // Azure Event Hubs is the feature being explained. The sample returns that focused object alongside a minimal run result.
        var endpoints = AzureEventHubsEndpoints();
        var result = LoadStrikeRunner
            .RegisterScenarios(Scenario("azure-event-hubs"))
            .WithRunnerKey(OrdersWorkflowPlaceholders.RunnerKey)
            .WithTestSuite("orders-reference")
            .WithTestName("azure-event-hubs")
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

            public static HttpEndpointDefinition HttpSource() =>

                new()

                {

                    Name = "orders-http-source",

                    Mode = TrafficEndpointMode.Produce,

                    TrackingField = TrackingFieldSelector.Parse("header:X-Correlation-Id"),

                    Url = $"{OrdersWorkflowPlaceholders.OrdersApiBaseUrl}/api/orders",

                    Method = "POST",

                    MessageHeaders = new Dictionary<string, string>(StringComparer.Ordinal)

                    {

                        ["X-Correlation-Id"] = OrdersWorkflowPlaceholders.ExampleOrderNumber,

                        ["X-Tenant"] = OrdersWorkflowPlaceholders.ExampleTenant,

                    },

                    MessagePayload = new

                    {

                        orderNumber = OrdersWorkflowPlaceholders.ExampleOrderNumber,

                        tenantId = OrdersWorkflowPlaceholders.ExampleTenant,

                    },

                    Auth = new HttpAuthOptions

                    {

                        Type = HttpAuthType.OAuth2ClientCredentials,

                        OAuth2ClientCredentials = new HttpOAuth2ClientCredentialsOptions

                        {

                            TokenEndpoint = $"{OrdersWorkflowPlaceholders.OrdersApiBaseUrl}/oauth/token",

                            ClientId = "dummy-client-id",

                            ClientSecret = "dummy-client-secret",

                            Scopes = ["orders.publish"],

                        },

                    },

                    TrackingPayloadSource = HttpTrackingPayloadSource.Request,

                    ResponseSource = HttpResponseSource.ResponseBody,

                };

            public static NatsEndpointDefinition NatsEndpoint() =>

                new()

                {

                    Name = "orders-nats",

                    Mode = TrafficEndpointMode.Consume,

                    TrackingField = TrackingFieldSelector.Parse("header:x-id"),

                    ServerUrl = OrdersWorkflowPlaceholders.NatsServerUrl,

                    Subject = "orders.events",

                    QueueGroup = "orders-workers",

                };

            public static object RabbitEndpoints() =>

                new

                {

                    Producer = new RabbitMqEndpointDefinition

                    {

                        Name = "orders-rabbit-producer",

                        Mode = TrafficEndpointMode.Produce,

                        TrackingField = TrackingFieldSelector.Parse("header:x-id"),

                        HostName = OrdersWorkflowPlaceholders.RabbitHost,

                        UserName = "guest",

                        Password = "guest",

                        QueueName = "orders.queue",

                        RoutingKey = "orders.route",

                        Durable = true,

                    },

                    Consumer = new RabbitMqEndpointDefinition

                    {

                        Name = "orders-rabbit-consumer",

                        Mode = TrafficEndpointMode.Consume,

                        TrackingField = TrackingFieldSelector.Parse("header:x-id"),

                        HostName = OrdersWorkflowPlaceholders.RabbitHost,

                        UserName = "guest",

                        Password = "guest",

                        QueueName = "orders.queue",

                        AutoAck = false,

                    },

                };

            public static object RedisEndpoints() =>

                new

                {

                    Producer = new RedisStreamsEndpointDefinition

                    {

                        Name = "orders-redis-producer",

                        Mode = TrafficEndpointMode.Produce,

                        TrackingField = TrackingFieldSelector.Parse("header:x-id"),

                        ConnectionString = OrdersWorkflowPlaceholders.RedisConnectionString,

                        StreamKey = "orders-stream",

                        MaxLength = 1000,

                    },

                    Consumer = new RedisStreamsEndpointDefinition

                    {

                        Name = "orders-redis-consumer",

                        Mode = TrafficEndpointMode.Consume,

                        TrackingField = TrackingFieldSelector.Parse("header:x-id"),

                        ConnectionString = OrdersWorkflowPlaceholders.RedisConnectionString,

                        StreamKey = "orders-stream",

                        ConsumerGroup = "orders-group",

                        ConsumerName = "orders-consumer",

                        ReadCount = 10,

                    },

                };

            public static object AzureEventHubsEndpoints() =>

                new

                {

                    Producer = new AzureEventHubsEndpointDefinition

                    {

                        Name = "orders-eventhub-producer",

                        Mode = TrafficEndpointMode.Produce,

                        TrackingField = TrackingFieldSelector.Parse("header:x-id"),

                        ConnectionString = OrdersWorkflowPlaceholders.AzureEventHubsConnectionString,

                        EventHubName = "orders",

                    },

                    Consumer = new AzureEventHubsEndpointDefinition

                    {

                        Name = "orders-eventhub-consumer",

                        Mode = TrafficEndpointMode.Consume,

                        TrackingField = TrackingFieldSelector.Parse("header:x-id"),

                        ConnectionString = OrdersWorkflowPlaceholders.AzureEventHubsConnectionString,

                        EventHubName = "orders",

                        ConsumerGroup = "$Default",

                        PartitionId = "0",

                        StartFromEarliest = true,

                    },

                };
}
