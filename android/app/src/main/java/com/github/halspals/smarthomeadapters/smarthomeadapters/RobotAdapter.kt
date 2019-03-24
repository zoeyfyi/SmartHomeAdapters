package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.BaseAdapter
import android.widget.ImageView
import android.widget.TextView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import org.jetbrains.anko.startActivity


class RobotAdapter (
        private val context: Context,
        private val robots: MutableList<Robot>,
        private val onClick: (Robot) -> Unit
) :  BaseAdapter() {

    init {
        // Add the dummy robot for starting the [RegisterRobotActivity]
        robots.add(Robot.ADD_ROBOT)
    }

    override fun getView(position: Int, convertView: View?, parent: ViewGroup?): View =
        if (position == count - 1) {
            assert(robots[position] == Robot.ADD_ROBOT) {
                "Getting view of last robot but it is not the ADD_ROBOT dummy; robot is " +
                        "${robots[position]}"
            }

            getAddRobotDummyView(convertView, parent)
        } else {
            getRealRobotView(position, convertView, parent)
    }

    private fun getRealRobotView(position: Int, convertView: View?, parent: ViewGroup?): View {

        // inflate card view
        val inflater = context.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val view: View = convertView ?: inflater.inflate(R.layout.view_robot_card, parent, false)

        // get internal views
        val robotNickname = view.findViewById<TextView>(R.id.robot_nickname_text_view)
        val robotCircle = view.findViewById<ImageView>(R.id.robot_circle_drawable)
        val robotIcon = view.findViewById<ImageView>(R.id.robot_image_view)

        // configure views
        robotCircle.setOnClickListener { onClick(robots[position]) }
        robotIcon.setImageResource(
            when (robots[position].robotType) {
                Robot.ROBOT_TYPE_SWITCH -> R.drawable.basic_lightbulb
                Robot.ROBOT_TYPE_THERMOSTAT -> R.drawable.basic_accelerator
                else -> R.drawable.basic_home
            }
        )
        robotNickname.text = robots[position].nickname

        return view
    }

    private fun getAddRobotDummyView(convertView: View?, parent: ViewGroup?): View {
        // inflate card view
        val inflater = context.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val view: View = convertView ?: inflater.inflate(R.layout.view_add_robot_card, parent, false)

        // get internal views
        val robotCircle = view.findViewById<ImageView>(R.id.robot_circle_drawable)

        // configure views
        robotCircle.setOnClickListener { context.startActivity<RegisterRobotActivity>() }

        return view
    }

    override fun getItem(position: Int): Any = robots[position]

    override fun getItemId(position: Int): Long = 0L

    override fun getCount(): Int = robots.size

}