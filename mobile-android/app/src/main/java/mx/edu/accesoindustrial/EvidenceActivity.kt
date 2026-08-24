package mx.edu.accesoindustrial

import android.Manifest
import android.graphics.Bitmap
import android.media.MediaRecorder
import android.os.Bundle
import android.view.Gravity
import android.widget.*
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.ContextCompat
import androidx.lifecycle.lifecycleScope
import com.google.gson.Gson
import kotlinx.coroutines.launch
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.RequestBody.Companion.asRequestBody
import okhttp3.RequestBody.Companion.toRequestBody
import java.io.File

class EvidenceActivity : AppCompatActivity() {
    private var photo: Bitmap? = null
    private var audioFile: File? = null
    private var recorder: MediaRecorder? = null
    private lateinit var photoPreview: ImageView
    private lateinit var audioButton: Button
    private lateinit var sendButton: Button
    private var recording = false
    private val picture = registerForActivityResult(ActivityResultContracts.TakePicturePreview()) { bitmap ->
        photo = bitmap
        if (bitmap != null) { photoPreview.setImageBitmap(bitmap); Toast.makeText(this, "Foto de evidencia agregada", Toast.LENGTH_SHORT).show() }
    }
    private val permissions = registerForActivityResult(ActivityResultContracts.RequestMultiplePermissions()) { }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        val root = LinearLayout(this).apply { orientation = LinearLayout.VERTICAL; setPadding(36, 36, 36, 36) }
        root.addView(TextView(this).apply { text = "Evidencia · ${MaintenanceDraft.machineName}"; textSize = 24f })
        root.addView(TextView(this).apply { text = "Agrega una fotografía obligatoria y una nota de voz opcional."; setPadding(0, 8, 0, 20) })
        photoPreview = ImageView(this).apply { adjustViewBounds = true; minimumHeight = 240; setImageResource(android.R.drawable.ic_menu_camera) }
        val photoButton = Button(this).apply { text = "Tomar foto de evidencia"; setOnClickListener { picture.launch(null) } }
        audioButton = Button(this).apply { text = "Grabar nota de voz (opcional)"; setOnClickListener { toggleRecording() } }
        sendButton = Button(this).apply { text = "Enviar mantenimiento"; setOnClickListener { submit() } }
        root.addView(photoPreview); root.addView(photoButton); root.addView(audioButton); root.addView(sendButton)
        setContentView(ScrollView(this).apply { addView(root) })
        permissions.launch(arrayOf(Manifest.permission.CAMERA, Manifest.permission.RECORD_AUDIO))
    }

    private fun toggleRecording() {
        if (recording) stopRecording() else startRecording()
    }
    private fun startRecording() {
        try {
            audioFile = File(cacheDir, "nota_${System.currentTimeMillis()}.m4a")
            recorder = MediaRecorder().apply { setAudioSource(MediaRecorder.AudioSource.MIC); setOutputFormat(MediaRecorder.OutputFormat.MPEG_4); setOutputFile(audioFile!!.absolutePath); setAudioEncoder(MediaRecorder.AudioEncoder.AAC); prepare(); start() }
            recording = true; audioButton.text = "Detener grabación"
        } catch (_: Exception) { Toast.makeText(this, "No se pudo iniciar la grabación", Toast.LENGTH_LONG).show() }
    }
    private fun stopRecording() {
        try { recorder?.stop() } catch (_: Exception) { audioFile = null }
        recorder?.release(); recorder = null; recording = false; audioButton.text = if (audioFile != null) "Nota de voz agregada · Grabar de nuevo" else "Grabar nota de voz (opcional)"
    }
    private fun submit() = lifecycleScope.launch {
        val image = photo ?: run { Toast.makeText(this@EvidenceActivity, "La foto de evidencia es obligatoria", Toast.LENGTH_LONG).show(); return@launch }
        sendButton.isEnabled = false; sendButton.text = "Enviando reporte..."
        try {
            val photoFile = File(cacheDir, "evidencia_${System.currentTimeMillis()}.jpg"); photoFile.outputStream().use { image.compress(Bitmap.CompressFormat.JPEG, 88, it) }
            val media = "text/plain".toMediaType()
            val details = MaintenanceDraft.checklist.map { MaintenanceDetail(it.itemId, it.complies) }
            val photoPart = MultipartBody.Part.createFormData("foto", photoFile.name, photoFile.asRequestBody("image/jpeg".toMediaType()))
            val audioPart = audioFile?.let { MultipartBody.Part.createFormData("audio", it.name, it.asRequestBody("audio/mp4".toMediaType())) }
            MaintenanceContainer.api.register(MaintenanceDraft.machineId.toRequestBody(media), Gson().toJson(details).toRequestBody(media), MaintenanceDraft.observations.toRequestBody(media), photoPart, audioPart)
            MaintenanceDraft.clear(); Toast.makeText(this@EvidenceActivity, "Mantenimiento enviado correctamente", Toast.LENGTH_LONG).show(); finish()
        } catch (_: Exception) { sendButton.isEnabled = true; sendButton.text = "Reintentar envío"; Toast.makeText(this@EvidenceActivity, "No se pudo enviar el reporte", Toast.LENGTH_LONG).show() }
    }
    override fun onDestroy() { if (recording) stopRecording(); super.onDestroy() }
}
