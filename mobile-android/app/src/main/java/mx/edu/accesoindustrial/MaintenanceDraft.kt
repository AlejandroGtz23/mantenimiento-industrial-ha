package mx.edu.accesoindustrial

data class DraftChecklistItem(val itemId: String, val complies: Boolean)
object MaintenanceDraft {
    var machineId: String = ""
    var machineName: String = ""
    var observations: String = ""
    var checklist: List<DraftChecklistItem> = emptyList()
    fun clear() { machineId = ""; machineName = ""; observations = ""; checklist = emptyList() }
}
