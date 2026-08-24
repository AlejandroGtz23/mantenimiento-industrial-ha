package mx.edu.accesoindustrial

import android.os.Bundle
import android.Manifest
import android.content.pm.PackageManager
import androidx.activity.result.contract.ActivityResultContracts
import androidx.core.content.ContextCompat
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import androidx.work.OneTimeWorkRequestBuilder
import androidx.work.WorkManager
import com.journeyapps.barcodescanner.ScanContract
import com.journeyapps.barcodescanner.ScanOptions
import kotlinx.coroutines.launch
import mx.edu.accesoindustrial.databinding.ActivityScannerBinding

class ScannerActivity : AppCompatActivity() { private lateinit var binding: ActivityScannerBinding
 private val scanner = registerForActivityResult(ScanContract()) { result -> result.contents?.let(::validate) }
 private val cameraPermission = registerForActivityResult(ActivityResultContracts.RequestPermission()) { granted -> if (granted) openScanner() else Toast.makeText(this, "Se requiere permiso de cámara", Toast.LENGTH_LONG).show() }
 override fun onCreate(savedInstanceState: Bundle?) { super.onCreate(savedInstanceState); binding = ActivityScannerBinding.inflate(layoutInflater); setContentView(binding.root); binding.scan.setOnClickListener { if (ContextCompat.checkSelfPermission(this, Manifest.permission.CAMERA) == PackageManager.PERMISSION_GRANTED) openScanner() else cameraPermission.launch(Manifest.permission.CAMERA) } }
 private fun openScanner() { scanner.launch(ScanOptions().setPrompt("Apunta al código QR").setBeepEnabled(true).setOrientationLocked(false)) }
 private fun validate(qr: String) { val device = "ANDROID-${android.os.Build.MODEL}"; lifecycleScope.launch { try { val response = AppContainer.api.validate(ValidateRequest(qr, device)); Toast.makeText(this@ScannerActivity, response.mensaje, Toast.LENGTH_LONG).show() } catch (_: Exception) { val database = LocalDatabase.get(this@ScannerActivity); database.pending().insert(PendingValidation(qr = qr, device = device)); WorkManager.getInstance(this@ScannerActivity).enqueue(OneTimeWorkRequestBuilder<SyncWorker>().build()); Toast.makeText(this@ScannerActivity, "Validación guardada para sincronizar", Toast.LENGTH_LONG).show() } } }
}
