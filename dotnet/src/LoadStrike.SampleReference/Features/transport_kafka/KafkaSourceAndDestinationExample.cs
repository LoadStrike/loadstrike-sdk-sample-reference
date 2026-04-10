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

namespace LoadStrike.SampleReference.Features.TransportKafka;

/// <summary>
/// /// This feature area focuses on Kafka producer and consumer setup, consumer groups, security, SASL, and OAuth bearer branches. ///
/// </summary>
public static class KafkaSourceAndDestinationExample
{
    public static object Build()
    {
        // Kafka Source And Destination is the feature being explained. The sample returns that focused object alongside a minimal run result.
        var endpoints = KafkaEndpoints();
        var result = LoadStrikeRunner
            .RegisterScenarios(Scenario("kafka-source-and-destination"))
            .WithRunnerKey(OrdersWorkflowPlaceholders.RunnerKey)
            .WithTestSuite("orders-reference")
            .WithTestName("kafka-source-and-destination")
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

            public static object KafkaEndpoints() =>

                new

                {

                    Producer = new KafkaEndpointDefinition

                    {

                        Name = "orders-kafka-producer",

                        Mode = TrafficEndpointMode.Produce,

                        TrackingField = TrackingFieldSelector.Parse("header:x-id"),

                        BootstrapServers = OrdersWorkflowPlaceholders.KafkaBootstrapServers,

                        Topic = OrdersWorkflowPlaceholders.OrderTopic,

                        SecurityProtocol = KafkaSecurityProtocolType.SaslSsl,

                        Sasl = new KafkaSaslOptions

                        {

                            Mechanism = KafkaSaslMechanismType.Plain,

                            Username = "dummy-user",

                            Password = "dummy-password",

                            AdditionalSettings = new Dictionary<string, string>(StringComparer.Ordinal)

                            {

                                ["socket.timeout.ms"] = "15000",

                            },

                        },

                    },

                    Consumer = new KafkaEndpointDefinition

                    {

                        Name = "orders-kafka-consumer",

                        Mode = TrafficEndpointMode.Consume,

                        TrackingField = TrackingFieldSelector.Parse("header:x-id"),

                        BootstrapServers = OrdersWorkflowPlaceholders.KafkaBootstrapServers,

                        Topic = OrdersWorkflowPlaceholders.OrderTopic,

                        ConsumerGroupId = "orders-sample-group",

                        StartFromEarliest = true,

                    },

                };

            public static KafkaEndpointDefinition KafkaOauthEndpoint() =>

                new()

                {

                    Name = "orders-kafka-oauth",

                    Mode = TrafficEndpointMode.Produce,

                    TrackingField = TrackingFieldSelector.Parse("header:x-id"),

                    BootstrapServers = OrdersWorkflowPlaceholders.KafkaBootstrapServers,

                    Topic = OrdersWorkflowPlaceholders.OrderTopic,

                    SecurityProtocol = KafkaSecurityProtocolType.SaslSsl,

                    Sasl = new KafkaSaslOptions

                    {

                        Mechanism = KafkaSaslMechanismType.OAuthBearer,

                        OAuthBearerTokenEndpointUrl =

                            $"{OrdersWorkflowPlaceholders.OrdersApiBaseUrl}/oauth/token",

                    },

                };
}
