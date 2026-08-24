package mx.edu.accesoindustrial

import android.content.Context
import androidx.room.*
import androidx.work.CoroutineWorker
import androidx.work.WorkerParameters
import com.google.gson.GsonBuilder
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

object AppContainer { val api: AccessApi by lazy { Retrofit.Builder().baseUrl("http://10.0.2.2:8081/").client(OkHttpClient()).addConverterFactory(GsonConverterFactory.create(GsonBuilder().create())).build().create(AccessApi::class.java) } }
@Entity(tableName = "pending_validations") data class PendingValidation(@PrimaryKey(autoGenerate = true) val id: Long = 0, val qr: String, val device: String)
@Dao interface PendingValidationDao { @Query("SELECT * FROM pending_validations") suspend fun all(): List<PendingValidation>; @Insert suspend fun insert(item: PendingValidation); @Delete suspend fun delete(item: PendingValidation) }
@Database(entities = [PendingValidation::class], version = 1) abstract class LocalDatabase : RoomDatabase() { abstract fun pending(): PendingValidationDao; companion object { fun get(context: Context) = Room.databaseBuilder(context, LocalDatabase::class.java, "access.db").build() } }
class SyncWorker(context: Context, params: WorkerParameters) : CoroutineWorker(context, params) { override suspend fun doWork(): Result { val dao = LocalDatabase.get(applicationContext).pending(); dao.all().forEach { item -> try { AppContainer.api.validate(ValidateRequest(item.qr, item.device)); dao.delete(item) } catch (_: Exception) { return Result.retry() } }; return Result.success() } }
