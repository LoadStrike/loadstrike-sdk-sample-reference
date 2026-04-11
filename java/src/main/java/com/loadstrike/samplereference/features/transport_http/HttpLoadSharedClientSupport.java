package com.loadstrike.samplereference.features.transport_http;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public final class HttpLoadSharedClientSupport {
  private static final ExecutorService SHARED_HTTP_EXECUTOR = Executors.newFixedThreadPool(64);
  private static final HttpClient SHARED_HTTP_CLIENT =
      HttpClient.newBuilder()
          .connectTimeout(Duration.ofSeconds(15))
          .executor(SHARED_HTTP_EXECUTOR)
          .version(HttpClient.Version.HTTP_1_1)
          .build();

  private HttpLoadSharedClientSupport() {}

  public static String postOrderStatusCode(String orderId, String tenantId) {
    String body =
        "{\"orderId\":\""
            + orderId
            + "\",\"tenantId\":\""
            + tenantId
            + "\",\"amount\":49.95}";

    HttpRequest request =
        HttpRequest.newBuilder(URI.create("https://api.example.com/orders"))
            .timeout(Duration.ofSeconds(15))
            .header("Content-Type", "application/json")
            .header("Accept-Encoding", "gzip, deflate")
            .POST(HttpRequest.BodyPublishers.ofString(body))
            .build();

    try {
      HttpResponse<String> response =
          SHARED_HTTP_CLIENT.send(request, HttpResponse.BodyHandlers.ofString());
      return Integer.toString(response.statusCode());
    } catch (IOException | InterruptedException exception) {
      if (exception instanceof InterruptedException) {
        Thread.currentThread().interrupt();
      }
      return "HTTP_EXCEPTION";
    }
  }

  public static void shutdownSharedHttpClient() {
    SHARED_HTTP_EXECUTOR.shutdownNow();
  }
}
