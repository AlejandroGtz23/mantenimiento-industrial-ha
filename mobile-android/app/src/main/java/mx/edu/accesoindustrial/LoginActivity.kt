package mx.edu.accesoindustrial

import android.content.Intent
import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import kotlinx.coroutines.launch
import mx.edu.accesoindustrial.databinding.ActivityLoginBinding
import retrofit2.HttpException
import java.io.IOException

class LoginActivity : AppCompatActivity() { private lateinit var binding: ActivityLoginBinding
 override fun onCreate(savedInstanceState: Bundle?) { super.onCreate(savedInstanceState); binding = ActivityLoginBinding.inflate(layoutInflater); setContentView(binding.root); binding.login.setOnClickListener { val number = binding.employeeNumber.text.toString().trim(); if (number.isBlank()) { binding.employeeNumber.error = "Ingresa tu número de empleado"; return@setOnClickListener }; lifecycleScope.launch { try { val login = MaintenanceContainer.api.login(MaintenanceLogin(number)); MaintenanceSession.token = login.token; startActivity(Intent(this@LoginActivity, MachinesActivity::class.java)); finish() } catch (error: HttpException) { binding.employeeNumber.error = if (error.code() == 401) "Técnico no encontrado o inactivo" else "Error del servidor (${error.code()})" } catch (_: IOException) { binding.employeeNumber.error = "No se pudo conectar al servidor" } catch (_: Exception) { binding.employeeNumber.error = "Error inesperado al iniciar sesión" } } } } }
