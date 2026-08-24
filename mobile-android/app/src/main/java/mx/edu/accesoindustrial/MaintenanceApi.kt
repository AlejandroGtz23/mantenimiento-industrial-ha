package mx.edu.accesoindustrial
import retrofit2.http.*
import okhttp3.MultipartBody
import okhttp3.RequestBody
data class MaintenanceLogin(val numero_empleado:String)
data class MaintenanceLoginResponse(val token:String,val tecnico:MaintenanceTechnician)
data class MaintenanceTechnician(val id:String,val nombre:String,val area:String)
data class MaintenanceMachine(val id:String,val codigo:String,val nombre:String,val ubicacion:String,val estado:String)
data class MaintenanceChecklistItem(val id:String,val descripcion:String,val orden:Int)
data class MaintenanceDetail(val item_checklist_id:String,val cumple:Boolean,val observacion_item:String="")
data class SubmittedMaintenance(val id:String,val estado:String)
interface MaintenanceApi { @POST("api/movil/login") suspend fun login(@Body value:MaintenanceLogin):MaintenanceLoginResponse; @GET("api/movil/maquinas") suspend fun machines():List<MaintenanceMachine>; @GET("api/movil/maquinas/{id}/checklist") suspend fun checklist(@Path("id") id:String):List<MaintenanceChecklistItem>; @Multipart @POST("api/movil/registrar-mantenimiento") suspend fun register(@Part("maquina_id") machine:RequestBody,@Part("detalles") details:RequestBody,@Part("observaciones") notes:RequestBody,@Part photo:MultipartBody.Part?=null,@Part audio:MultipartBody.Part?=null):SubmittedMaintenance }
