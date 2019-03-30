package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.BaseAdapter
import android.widget.TextView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.UseCase

class UseCaseAdapter(private val context: Context, private val useCases: List<UseCase>)
    : BaseAdapter() {

    override fun getView(position: Int, convertView: View?, parent: ViewGroup): View {

        val inflater = context.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val view = convertView ?: inflater.inflate(R.layout.attachment_list_item, parent, false)
        view.findViewById<TextView>(R.id.list_item_header).text = useCases[position].name
        view.findViewById<TextView>(R.id.list_item_description).text = useCases[position].name // TODO change obvs
        return view
    }

    override fun getItem(pos: Int): Any = useCases[pos]
    override fun getCount(): Int = useCases.size
    override fun getItemId(pos: Int): Long = 0L
}