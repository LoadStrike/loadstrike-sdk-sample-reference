package com.loadstrike.samplereference.features.transport_http;

import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeResponse;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeRunner;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeScenario;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeSimulation;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeStep;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.UUID;

public final class HttpLoadSharedClientExample {
  private HttpLoadSharedClientExample() {}

  public static Object build() {
    try {
      var result =
          LoadStrikeRunner.registerScenarios(scenario())
              .withRunnerKey("runner_dummy_orders_reference")
              .withTestSuite("orders-reference")
              .withTestName("http-load-shared-client")
              .withoutReports()
              .run();

      return mapOf(
          "note", "Shut down the shared HTTP client executor when the run ends.", "Result", result);
    } finally {
      HttpLoadSharedClientSupport.shutdownSharedHttpClient();
    }
  }

  public static LoadStrikeScenario scenario() {
    return LoadStrikeScenario.create(
            "http-load-shared-client",
            context -> {
              var orderId = "ord-" + UUID.randomUUID();
              var tenantId = "tenant-a";
              return LoadStrikeStep.run(
                      "POST /orders",
                      context,
                      () -> {
                        var statusCode =
                            HttpLoadSharedClientSupport.postOrderStatusCode(orderId, tenantId);
                        return statusCode.startsWith("2")
                            ? LoadStrikeResponse.ok(statusCode)
                            : LoadStrikeResponse.fail(statusCode);
                      })
                  .asReply();
            })
        .withLoadSimulations(LoadStrikeSimulation.inject(25, 1d, 30d));
  }

  public static Map<String, Object> mapOf(Object... values) {
    var result = new LinkedHashMap<String, Object>();
    for (int index = 0; index < values.length; index += 2) {
      result.put(String.valueOf(values[index]), values[index + 1]);
    }
    return result;
  }
}
