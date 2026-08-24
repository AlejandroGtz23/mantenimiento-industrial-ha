package mx.edu.accesoindustrial
import android.content.Intent
import android.os.Bundle
import android.widget.*
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import kotlinx.coroutines.launch
class MachinesActivity:AppCompatActivity(){override fun onCreate(savedInstanceState:Bundle?){super.onCreate(savedInstanceState);val list=ListView(this);setContentView(list);lifecycleScope.launch{try{val machines=MaintenanceContainer.api.machines();list.adapter=ArrayAdapter(this@MachinesActivity,android.R.layout.simple_list_item_2,android.R.id.text1,machines.map{"${it.codigo} · ${it.nombre}\n${it.ubicacion} — ${it.estado}"});list.setOnItemClickListener{_,_,position,_->startActivity(Intent(this@MachinesActivity,ChecklistActivity::class.java).putExtra("machine_id",machines[position].id).putExtra("machine_name",machines[position].nombre))}}catch(e:Exception){Toast.makeText(this@MachinesActivity,"No se pudieron cargar las máquinas",Toast.LENGTH_LONG).show()}}}}
