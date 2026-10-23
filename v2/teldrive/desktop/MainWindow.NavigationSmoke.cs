using System.IO;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Media;
using System.Windows.Media.Imaging;

namespace Teldrive.Desktop;

public partial class MainWindow
{
    internal async Task VerifyNavigationAsync()
    {
        await ShowDrive(); await Task.Delay(1200);
        Logs.AppendText("Prueba de actividad: registro visible sin una capa del navegador.\n");
        LogsClick(this, new RoutedEventArgs()); await Task.Delay(200); UpdateLayout();
        if (!LogsView.IsVisible || Browser.IsVisible || SetupView.IsVisible || AdvancedView.IsVisible) throw new IOException("Actividad no ocultó las demás vistas.");
        var image = new RenderTargetBitmap((int)LogsView.ActualWidth, (int)LogsView.ActualHeight, 96, 96, PixelFormats.Pbgra32); image.Render(LogsView);
        // Log text can cross the center pixel. Sample a patch so a light glyph
        // is not mistaken for a light browser surface covering the activity view.
        const int sampleSize = 20;
        var pixels = new byte[sampleSize * sampleSize * 4];
        image.CopyPixels(new Int32Rect((int)LogsView.ActualWidth / 2 - sampleSize / 2, (int)LogsView.ActualHeight / 2 - sampleSize / 2, sampleSize, sampleSize), pixels, sampleSize * 4, 0);
        var darkPixels = Enumerable.Range(0, sampleSize * sampleSize).Count(i => pixels[i * 4 + 3] > 0 && pixels[i * 4] <= 100 && pixels[i * 4 + 1] <= 100 && pixels[i * 4 + 2] <= 100);
        var encoder = new PngBitmapEncoder(); encoder.Frames.Add(BitmapFrame.Create(image)); using (var output = File.Create(Path.Combine(runtime.DataRoot, "activity.png"))) encoder.Save(output);
        await File.WriteAllTextAsync(Path.Combine(runtime.DataRoot, "activity-layout.json"), System.Text.Json.JsonSerializer.Serialize(new { state = WindowState.ToString(), width = LogsView.ActualWidth, height = LogsView.ActualHeight, logWidth = Logs.ActualWidth, logHeight = Logs.ActualHeight, logBackground = Logs.Background.ToString(), darkPixels, sampledPixels = sampleSize * sampleSize }));
        if (darkPixels < sampleSize * sampleSize * 0.6) throw new IOException("Actividad muestra una capa clara en el centro del registro.");
        // Exercise the actual header button after initialization, as a user would.
        busy = false; DriveButton.RaiseEvent(new RoutedEventArgs(Button.ClickEvent));
        for (var attempt=0;attempt<40&&!Browser.IsVisible;attempt++) await Task.Delay(50);
        if (!Browser.IsVisible || LogsView.IsVisible) throw new IOException("El botón Teldrive no recuperó el Drive.");
        await Task.Delay(100); busy = true;
        ShowAdvanced(); if (!AdvancedView.IsVisible || Browser.IsVisible) throw new IOException("Rclone no ocultó el navegador.");
        SetupClick(this,new RoutedEventArgs()); if (!SetupView.IsVisible || AdvancedView.IsVisible) throw new IOException("Configuración no recuperó su vista.");
        await ShowDrive();
    }
}
