package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.*
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.UseCase

class UseCaseAdapter(context: Context, private val useCases: List<UseCase>)
    : ArrayAdapter<UseCase>(context, 0, useCases) {

    internal var selectedUseCasePos: Int? = null

    private val tag = "UseCaseAdapter"

    override fun getView(position: Int, convertView: View?, parent: ViewGroup): View {
        val inflater = context.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val view = convertView ?: inflater.inflate(R.layout.attachment_list_item, parent, false)
        view.findViewById<TextView>(R.id.list_item_header).text = useCases[position].name.capitalize()
        view.findViewById<TextView>(R.id.list_item_description).text = useCases[position].description

        val checkBox = view.findViewById<RadioButton>(R.id.radio)
        checkBox.isChecked = position == selectedUseCasePos
        Log.v(tag, "Item at $position is selected: ${position == selectedUseCasePos}")

        return view
    }
}