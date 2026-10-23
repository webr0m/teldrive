using System.IO;
using System.Net.Http;
using System.Reflection;
using System.Text.Json;
using System.Windows;
using System.Windows.Threading;

namespace Teldrive.Desktop;

public partial class MainWindow
{
    private async Task VerifyRemoteConnectionAsync()
    {
        var original = await File.ReadAllTextAsync(runtime.ConfigPath);
        var endpoint = runtime.ServerUri ?? throw new IOException("Remote TOML did not initialize Desktop without desktop.json.");
        void NoLocalServices() {
            if (!runtime.Remote || Directory.Exists(Path.Combine(runtime.DataRoot, "database")) || new[] { "server", "postgres" }.Any(name => typeof(RuntimeSession).GetField(name, BindingFlags.Instance | BindingFlags.NonPublic)!.GetValue(runtime) is not null)) throw new IOException("Remote mode started local services.");
        }
        NoLocalServices();
        RefreshConfiguration();
        if (LocalServerSetup.IsEnabled) throw new IOException("The local setup form can overwrite a remote connection.");
        await runtime.StartAsync(); await ShowDrive();
        var rendered = false;
        for (var attempt = 0; attempt < 80; attempt++) {
            var value = await Browser.CoreWebView2.ExecuteScriptAsync("JSON.stringify({title:document.title,text:document.body.innerText})");
            var text = JsonSerializer.Deserialize<string>(value) ?? "";
            if (text.Contains("Teldrive", StringComparison.OrdinalIgnoreCase) && text.Contains("Telegram", StringComparison.OrdinalIgnoreCase)) { rendered = true; break; }
            await Task.Delay(250);
        }
        if (!rendered || Browser.Source.GetLeftPart(UriPartial.Authority) != endpoint.GetLeftPart(UriPartial.Authority)) throw new IOException("Remote login UI did not render on the configured origin.");
        await Browser.CoreWebView2.CallDevToolsProtocolMethodAsync("Runtime.evaluate", JsonSerializer.Serialize(new { expression = "document.documentElement.requestFullscreen()", userGesture = true, awaitPromise = true }));
        await Dispatcher.Yield(DispatcherPriority.ApplicationIdle);
        if (!browserFullscreen.Active || AppHeader.Visibility != Visibility.Collapsed) throw new IOException("Remote web fullscreen failed.");
        await Browser.CoreWebView2.ExecuteScriptAsync("document.exitFullscreen()");
        await Dispatcher.Yield(DispatcherPriority.ApplicationIdle);
        if (browserFullscreen.Active) throw new IOException("Remote fullscreen did not restore the window.");
        var rejected = false;
        try { await runtime.SaveConfigAsync("[desktop]\nserver-url = \"http://remote.example.com\"\n"); } catch (IOException) { rejected = true; }
        if (!rejected || await File.ReadAllTextAsync(runtime.ConfigPath) != original || runtime.ServerUri != endpoint) throw new IOException("Invalid remote config replaced the working configuration.");
        NoLocalServices();
        await runtime.StopAsync();
        if (runtime.Running) throw new IOException("Remote disconnect did not clear the connection state.");
        await runtime.SaveConfigAsync("[desktop]\nserver-url = \"http://127.0.0.1:1\"\n");
        var unavailableRejected = false;
        try { await runtime.StartAsync(); } catch (HttpRequestException) { unavailableRejected = true; } catch (IOException) { unavailableRejected = true; }
        if (!unavailableRejected || runtime.Running) throw new IOException("An unavailable remote server was treated as connected.");
        NoLocalServices();
        await runtime.SaveConfigAsync(original); await runtime.StartAsync();
        NoLocalServices();
        await runtime.StopAsync();
        using var http = new HttpClient(new HttpClientHandler { AllowAutoRedirect = false }) { Timeout = TimeSpan.FromSeconds(25) };
        using var health = await http.GetAsync(new Uri(endpoint, "/api/health/ready"));
        if (!health.IsSuccessStatusCode) throw new IOException("The central server stopped after disconnecting Desktop.");
        await File.WriteAllTextAsync(Path.Combine(runtime.DataRoot, "remote-smoke-test.json"), JsonSerializer.Serialize(new { success = true, endpoint = endpoint.GetLeftPart(UriPartial.Authority), configOnlyStartup = true, loginUi = rendered, browserFullscreen = true, noLocalServices = true, rejectedUnsafeUrl = rejected, rejectedUnavailableServer = unavailableRejected, configurationPreserved = await File.ReadAllTextAsync(runtime.ConfigPath) == original, centralServerStillReady = true, stopped = !runtime.Running }));
    }
}
