package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.GridView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot

class RobotsFragment : Fragment() {

    private val fTag = "RobotFragment"

    private lateinit var robotGrid: GridView

    // TODO: remove once we have a real data source
    private val robotIcons = listOf(
        R.drawable.basic_accelerator,
        R.drawable.basic_chronometer,
        R.drawable.basic_home,
        R.drawable.basic_key,
        R.drawable.basic_lightbulb,
        R.drawable.basic_lock,
        R.drawable.basic_lock_open
    )

    // TODO: get list of robots from REST API
    private val robots = (1..20).map {
        Robot("Robot $it", robotIcons[it % robotIcons.size])
    }

    override fun onCreateView(
            inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_robots, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        robotGrid = view.findViewById(R.id.RobotGrid)
        robotGrid.adapter = RobotAdapter(view.context, robots) { robot ->
            Log.d(fTag, "Clicked robot: \"${robot.nickname}\"")
            (activity as MainActivity).startFragment(RobotFragment(), true)
        }
    }
}
