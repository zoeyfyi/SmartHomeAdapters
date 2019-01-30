package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.view.View
import android.view.ViewGroup
import android.widget.BaseAdapter
import android.widget.TextView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot

class RobotAdapter (private val context: Context) :  BaseAdapter() {

    private val robots = arrayOf(
        Robot("Living room light"),
        Robot("Garage fridge"),
        Robot("Garden gate"),
        Robot("Thermostat")
    )

    override fun getView(position: Int, convertView: View?, parent: ViewGroup?): View {
        val textView: TextView
        if (convertView == null) {
            // if it's not recycled, initialize some attributes
            textView = TextView(context)
            textView.setPadding(8, 8, 8, 8)
        } else {
            textView = convertView as TextView
        }

        textView.text = robots[position].nickname
        return textView
    }

    override fun getItem(position: Int): Any = robots[position]

    override fun getItemId(position: Int): Long = 0L

    override fun getCount(): Int = robots.size

}