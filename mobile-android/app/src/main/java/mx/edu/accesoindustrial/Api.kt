package mx.edu.accesoindustrial

import okhttp3.MultipartBody
import okhttp3.RequestBody
import retrofit2.http.*

data class LoginRequest(val numero_empleado: String, val rol: String)
data class LoginResponse(val numero_empleado: String, val rol: String, val color_fondo: String)
data class SelfieResponse(val sesion_id: String, val codigo_qr: String, val expira_en: String)
data class ValidateRequest(val codigo_qr: String, val dispositivo_oficial: String, val rostro_coincide: Boolean = true)
data class ValidateResponse(val resultado: String, val mensaje: String)
interface AccessApi { @POST("api/movil/iniciar-sesion") suspend fun login(@Body request: LoginRequest): LoginResponse; @Multipart @POST("api/movil/registrar-selfie") suspend fun selfie(@Part selfie: MultipartBody.Part, @Part("numero_empleado") number: RequestBody, @Part("color_fondo") color: RequestBody): SelfieResponse; @POST("api/movil/validar-qr") suspend fun validate(@Body request: ValidateRequest): ValidateResponse }
