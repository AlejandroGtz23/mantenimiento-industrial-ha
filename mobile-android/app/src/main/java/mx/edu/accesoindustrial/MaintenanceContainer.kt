package mx.edu.accesoindustrial
import okhttp3.Interceptor
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
object MaintenanceSession { var token="" }
// La app entra por Nginx en la VM, no por una réplica del backend directamente.
object MaintenanceContainer { private val client=OkHttpClient.Builder().addInterceptor(Interceptor{chain->chain.proceed(chain.request().newBuilder().apply{if(MaintenanceSession.token.isNotBlank())header("Authorization","Bearer ${MaintenanceSession.token}")}.build())}).build(); val api:MaintenanceApi by lazy{Retrofit.Builder().baseUrl("http://192.168.1.48/").client(client).addConverterFactory(GsonConverterFactory.create()).build().create(MaintenanceApi::class.java)} }
