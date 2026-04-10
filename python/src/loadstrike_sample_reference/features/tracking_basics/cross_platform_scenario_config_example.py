from __future__ import annotations

from typing import Any

from playwright.async_api import async_playwright
from selenium.webdriver import ChromeOptions

from loadstrike_sdk import (
    AzureEventHubsEndpointDefinition,
    CrossPlatformScenarioConfigurator,
    DatadogReportingSink,
    DelegateStreamEndpointDefinition,
    GrafanaLokiReportingSink,
    HttpAuthOptions,
    HttpEndpointDefinition,
    HttpOAuth2ClientCredentialsOptions,
    InfluxDbReportingSink,
    KafkaEndpointDefinition,
    KafkaSaslOptions,
    LoadStrikeMetric,
    LoadStrikePluginData,
    LoadStrikePluginDataTable,
    LoadStrikeReportFormat,
    LoadStrikeResponse,
    LoadStrikeRunner,
    LoadStrikeScenario,
    LoadStrikeSimulation,
    LoadStrikeStep,
    LoadStrikeThreshold,
    NatsEndpointDefinition,
    OtelCollectorReportingSink,
    PushDiffusionEndpointDefinition,
    RabbitMqEndpointDefinition,
    RedisStreamsEndpointDefinition,
    SplunkReportingSink,
    TimescaleDbReportingSink,
    TrackingFieldSelector,
    TrackingPayloadBuilder,
)

FEATURE_OVERVIEW = "This feature area focuses on selectors, run modes, and cross-platform tracking configuration needed to correlate generated traffic."


def build():
    # Cross Platform Scenario Config keeps the scenario setup focused on one behavior. Change the option below first when comparing suboptions.
    scenario = _tracked_scenario()
    return (LoadStrikeRunner.register_scenarios(scenario).with_runner_key(placeholders.RUNNER_KEY).with_test_suite("orders-reference").with_test_name("cross-platform-scenario-config").without_reports().run())


class placeholders:
    RUNNER_KEY = "runner_dummy_orders_reference"
    ORDERS_API_BASE_URL = "https://orders.example.test"
    KAFKA_BOOTSTRAP_SERVERS = "localhost:9092"
    NATS_SERVER_URL = "nats://localhost:4222"
    REDIS_CONNECTION_STRING = "localhost:6379,abortConnect=false"
    RABBIT_HOST = "localhost"
    AZURE_EVENT_HUBS_CONNECTION_STRING = (
        "Endpoint=sb://orders.example.test/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=dummy"
    )
    REPORT_FOLDER = "./artifacts/reports"
    ORDER_TOPIC = "orders.created"
    EXAMPLE_ORDER_NUMBER = "ORD-10001"
    EXAMPLE_TENANT = "demo-tenant"

    @staticmethod
    def scenario_name(suffix: str) -> str:
        return f"orders.{suffix}"


def _step_reply(context: Any):
    return LoadStrikeStep.run(
        "publish-order",
        context,
        lambda: LoadStrikeResponse.ok(
            {
                "orderNumber": placeholders.EXAMPLE_ORDER_NUMBER,
                "tenant": placeholders.EXAMPLE_TENANT,
            },
            "201",
            128,
            "created",
            4.5,
        ),
    ).as_reply()


def _scenario(name_suffix: str = "publish"):
    return (
        LoadStrikeScenario.create(placeholders.scenario_name(name_suffix), _step_reply)
        .with_init(lambda _context: None)
        .with_clean(lambda _context: None)
        .with_weight(2)
        .without_warm_up()
        .with_max_fail_count(3)
        .with_restart_iteration_on_fail(True)
    )


def _http_source():
    return HttpEndpointDefinition(
        name="orders-http-source",
        mode="Produce",
        tracking_field=TrackingFieldSelector.parse("header:X-Correlation-Id"),
        url=f"{placeholders.ORDERS_API_BASE_URL}/api/orders",
        method="POST",
        message_headers={"X-Tenant": placeholders.EXAMPLE_TENANT},
        message_payload={
            "orderNumber": placeholders.EXAMPLE_ORDER_NUMBER,
            "tenantId": placeholders.EXAMPLE_TENANT,
        },
        auth=HttpAuthOptions(
            type="OAuth2ClientCredentials",
            oauth2_client_credentials=HttpOAuth2ClientCredentialsOptions(
                token_url=f"{placeholders.ORDERS_API_BASE_URL}/oauth/token",
                client_id="dummy-client-id",
                client_secret="dummy-client-secret",
                scope="orders.publish",
            ),
        ),
        tracking_payload_source="Request",
        response_source="Body",
    )


def _http_destination():
    return HttpEndpointDefinition(
        name="orders-http-destination",
        mode="Consume",
        tracking_field=TrackingFieldSelector.parse("json:$.trackingId"),
        gather_by_field=TrackingFieldSelector.parse("json:$.tenantId"),
        url=f"{placeholders.ORDERS_API_BASE_URL}/api/order-events",
        method="GET",
        response_source="Body",
        consume_json_array_response=True,
        consume_array_path="$.items",
    )


def _tracking_configuration():
    return {
        "Source": _http_source(),
        "Destination": _http_destination(),
        "RunMode": "GenerateAndCorrelate",
        "CorrelationTimeoutSeconds": 30,
        "TimeoutSweepIntervalSeconds": 1,
        "TimeoutBatchSize": 200,
        "TimeoutCountsAsFailure": True,
        "MetricPrefix": "orders_tracking",
        "ExecuteOriginalScenarioRun": False,
    }


def _tracked_scenario():
    return CrossPlatformScenarioConfigurator.configure(
        _scenario("tracking"),
        _tracking_configuration(),
    )
