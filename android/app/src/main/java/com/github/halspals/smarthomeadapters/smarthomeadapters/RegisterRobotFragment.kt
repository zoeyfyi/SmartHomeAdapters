package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import kotlinx.android.synthetic.main.fragment_register_robot.*


class RegisterRobotFragment : Fragment() {

    private var robotId: String? = null

    private val fTag = "RegisterRobotFragment"

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        robotId = arguments?.getString("robotId")

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_register_robot, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        register_button.setOnClickListener { _ ->
            if (registerRobot()) {
                (activity as MainActivity).startFragment(RobotsFragment(), true)
            }
        }
    }

    private fun registerRobot(): Boolean {
        // TODO read the inputs, send a registration query to the server,
        // parse the result
        return true
    }
}
