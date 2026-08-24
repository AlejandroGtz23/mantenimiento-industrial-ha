package mx.edu.accesoindustrial

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.launch

class AccessViewModel : ViewModel() {
    private val _message = MutableLiveData<String>()
    val message: LiveData<String> = _message

    fun validate(qr: String, device: String, onOffline: suspend () -> Unit) = viewModelScope.launch {
        try { _message.postValue(AppContainer.api.validate(ValidateRequest(qr, device)).mensaje) }
        catch (_: Exception) { onOffline(); _message.postValue("Validación pendiente de sincronización") }
    }
}
