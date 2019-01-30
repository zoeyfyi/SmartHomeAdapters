package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.BaseAdapter
import android.widget.ImageView
import android.widget.TextView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot


class RobotAdapter (private val context: Context, private val robots: List<Robot>, private val onClick: (Robot) -> Unit) :  BaseAdapter() {

    override fun getView(position: Int, convertView: View?, parent: ViewGroup?): View {
        // inflate card view
        val inflater = context.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val view: View = convertView ?: inflater.inflate(R.layout.view_robot_card, parent, false)

        // get internal views
        val robotIcon = view.findViewById<ImageView>(R.id.robot_icon_image_view)
        val robotNickname = view.findViewById<TextView>(R.id.robot_nickname_text_view)

        // configure views
        view.setOnClickListener { onClick(robots[position]) }
        robotIcon.setImageResource(robots[position].iconDrawable)
        robotNickname.text = robots[position].nickname

        return view
    }

    override fun getItem(position: Int): Any = robots[position]

    override fun getItemId(position: Int): Long = 0L

    override fun getCount(): Int = robots.size

}