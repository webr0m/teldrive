using System.Diagnostics;
using System.IO;
using System.Text.Json;
using System.Windows;
using System.Windows.Threading;

namespace Teldrive.Desktop;

public partial class MainWindow
{
    internal async Task VerifyWindowReopenAsync()
    {
        testingTray = true;
        preferences = new() { CloseToTray = true, MinimizeToTray = true };
        tray = new DesktopTray(RestoreFromTray, () => { RestoreFromTray(); ShowAdvanced(); }, () => { }, ExitDesktop);
        var cases = new List<string>();
        try {
            foreach (var state in new[] { WindowState.Maximized, WindowState.Normal }) {
                foreach (var route in new[] { "direct", "tray-menu", "second-launch" }) {
                    WindowState = state;
                    // Starting in the tray suppresses activation. Reopening must
                    // handle that flag even after the user maximizes the window.
                    ShowActivated = false;
                    Close();
                    if (IsVisible || !tray.Visible || closing) throw new IOException("Close did not keep the application in the tray.");
                    if (route == "direct") RestoreFromTray();
                    else if (route == "tray-menu") tray.InvokeMenuForSmoke(0);
                    else {
                        var start = new ProcessStartInfo(Environment.ProcessPath!) { UseShellExecute = false, CreateNoWindow = true };
                        start.ArgumentList.Add("--data-root"); start.ArgumentList.Add(runtime.DataRoot);
                        using var child = Process.Start(start) ?? throw new IOException("Could not relaunch the application.");
                        await child.WaitForExitAsync().WaitAsync(TimeSpan.FromSeconds(20));
                        if (child.ExitCode != 0) throw new IOException("The second application instance failed.");
                    }
                    for (var attempt = 0; attempt < 100 && !IsVisible; attempt++) await Task.Delay(50);
                    await Dispatcher.Yield(DispatcherPriority.ApplicationIdle);
                    if (!IsVisible || WindowState != state || !ShowActivated) throw new IOException($"Reopen failed for {state} through {route}.");
                    cases.Add($"{state}:{route}");
                }
            }
            WindowState = WindowState.Minimized;
            if (IsVisible || !tray.Visible) throw new IOException("Minimize did not hide the application.");
            RestoreFromTray();
            if (!IsVisible || WindowState != WindowState.Normal) throw new IOException("The minimized window did not reopen normally.");
            cases.Add("Minimized:restore");
            await File.WriteAllTextAsync(Path.Combine(runtime.DataRoot, "window-reopen-smoke-test.json"), JsonSerializer.Serialize(new { success = true, cases, configurationWritten = false, servicesStarted = false }));
        } finally { preferences = new(); tray.Dispose(); tray = null; testingTray = false; }
    }
}
