package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import kotlinx.android.synthetic.main.fragment_edit_robot.*
import kotlinx.android.synthetic.main.fragment_nickname.*
import kotlinx.android.synthetic.main.view_robot_card.*
import org.jetbrains.anko.clearTask
import org.jetbrains.anko.intentFor

class EditRobotFragment : Fragment() {

    private val fTag = "QRFragment"

    private val parent by lazy { activity as MainActivity }

    override fun onCreateView(
            inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?
    ): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_nickname, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        finish_edit_image_view.setOnClickListener { _ -> parent.startFragment(RobotsFragment()) }

        recalibrate_layout.setOnClickListener { _ -> /*TODO*/ }
        rename_layout.setOnClickListener { _ -> /*TODO*/ }
        delete_layout.setOnClickListener { _ -> /*TODO*/ }

    }

    internal fun setRobotView(robot: Robot) {
        // inflate card view
        val inflater = parent.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val robotView = inflater.inflate(R.layout.view_robot_card, robot_layout.parent as ViewGroup, false)

        // get internal views
        val robotNickname = robotView.findViewById<TextView>(R.id.robot_nickname_text_view)
        val robotIcon = robotView.findViewById<ImageView>(R.id.robot_image_view)


        // Set the icon accordingly
        when (robot.robotType) {
            Robot.ROBOT_TYPE_SWITCH -> {
                robotIcon.setImageResource(R.drawable.basic_lightbulb)
            }

            Robot.ROBOT_TYPE_THERMOSTAT -> {
                robotIcon.setImageResource(R.drawable.basic_accelerator)
            }

            else -> TODO("NO OTHER ROBOT TYPE EXPECTED")

        }
        robotNickname.text = robot.nickname

        robot_layout.addView(robotView)
    }
}
