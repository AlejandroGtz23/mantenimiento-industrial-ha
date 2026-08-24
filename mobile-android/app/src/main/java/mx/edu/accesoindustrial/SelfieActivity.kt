package mx.edu.accesoindustrial

import android.content.Intent
import android.Manifest
import android.content.pm.PackageManager
import android.os.Bundle
import java.io.File
import androidx.appcompat.app.AppCompatActivity
import androidx.camera.core.*
import androidx.camera.lifecycle.ProcessCameraProvider
import androidx.core.content.ContextCompat
import androidx.activity.result.contract.ActivityResultContracts
import mx.edu.accesoindustrial.databinding.ActivitySelfieBinding

class SelfieActivity : AppCompatActivity() { private lateinit var binding: ActivitySelfieBinding; private var capture: ImageCapture? = null
 private val cameraPermission = registerForActivityResult(ActivityResultContracts.RequestPermission()) { granted -> if (granted) startCamera() else binding.colorHint.text = "Se requiere permiso de cámara para tomar la selfie" }
 override fun onCreate(savedInstanceState: Bundle?) { super.onCreate(savedInstanceState); binding = ActivitySelfieBinding.inflate(layoutInflater); setContentView(binding.root); binding.colorHint.text = "Fondo requerido: ${getSharedPreferences("session", MODE_PRIVATE).getString("color", "")}"; if (ContextCompat.checkSelfPermission(this, Manifest.permission.CAMERA) == PackageManager.PERMISSION_GRANTED) startCamera() else cameraPermission.launch(Manifest.permission.CAMERA); binding.takePhoto.setOnClickListener { takePhoto() } }
 private fun startCamera() { val future = ProcessCameraProvider.getInstance(this); future.addListener({ val provider = future.get(); capture = ImageCapture.Builder().build(); provider.unbindAll(); provider.bindToLifecycle(this, CameraSelector.DEFAULT_FRONT_CAMERA, Preview.Builder().build().also { it.surfaceProvider = binding.preview.surfaceProvider }, capture) }, ContextCompat.getMainExecutor(this)) }
 private fun takePhoto() { val file = File.createTempFile("selfie", ".jpg", cacheDir); capture?.takePicture(ImageCapture.OutputFileOptions.Builder(file).build(), ContextCompat.getMainExecutor(this), object: ImageCapture.OnImageSavedCallback { override fun onError(exception: ImageCaptureException) {} ; override fun onImageSaved(result: ImageCapture.OutputFileResults) { startActivity(Intent(this@SelfieActivity, QrActivity::class.java).putExtra("file", file.absolutePath)) } }) }
}
