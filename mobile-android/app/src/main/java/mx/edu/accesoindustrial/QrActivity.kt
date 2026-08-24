package mx.edu.accesoindustrial

import android.graphics.Bitmap
import android.os.Bundle
import android.os.CountDownTimer
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import com.google.zxing.BarcodeFormat
import com.google.zxing.MultiFormatWriter
import com.google.zxing.common.BitMatrix
import kotlinx.coroutines.launch
import mx.edu.accesoindustrial.databinding.ActivityQrBinding
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.RequestBody.Companion.asRequestBody
import okhttp3.RequestBody.Companion.toRequestBody
import java.io.File

class QrActivity : AppCompatActivity() { private lateinit var binding: ActivityQrBinding
 override fun onCreate(savedInstanceState: Bundle?) { super.onCreate(savedInstanceState); binding = ActivityQrBinding.inflate(layoutInflater); setContentView(binding.root); val file = File(intent.getStringExtra("file") ?: return); val prefs = getSharedPreferences("session", MODE_PRIVATE); lifecycleScope.launch { try { val response = AppContainer.api.selfie(MultipartBody.Part.createFormData("selfie", file.name, file.asRequestBody("image/jpeg".toMediaType())), prefs.getString("employee", "")!!.toRequestBody("text/plain".toMediaType()), prefs.getString("color", "")!!.toRequestBody("text/plain".toMediaType())); binding.qrImage.setImageBitmap(toBitmap(MultiFormatWriter().encode(response.codigo_qr, BarcodeFormat.QR_CODE, 600, 600))); object : CountDownTimer(120000, 1000) { override fun onTick(ms: Long) { binding.timer.text = "Expira en ${ms / 1000}s" }; override fun onFinish() { binding.timer.text = "QR expirado" } }.start() } catch (_: Exception) { binding.timer.text = "No se pudo registrar la selfie" } } }
 private fun toBitmap(matrix: BitMatrix): Bitmap { val bitmap = Bitmap.createBitmap(matrix.width, matrix.height, Bitmap.Config.RGB_565); for (x in 0 until matrix.width) for (y in 0 until matrix.height) bitmap.setPixel(x, y, if (matrix[x,y]) android.graphics.Color.BLACK else android.graphics.Color.WHITE); return bitmap }
}
