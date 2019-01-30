package com.github.halspals.smarthomeadapters.smarthomeadapters


import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import android.widget.Toast
import org.jetbrains.anko.toast

class RobotFragment : Fragment() {

    lateinit var robotId: String

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {

        // get robotId from bundle
        val robotIdArgument = arguments?.getString("robotId")
        if (robotIdArgument == null) {
            // no id passed, try to go back
            Log.d(tag, "No robotId passed to robotFragment")
            context?.toast("Oops, something went wrong")
            fragmentManager?.popBackStack()
            return null
        }
        robotId = robotIdArgument

        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_robot, container, false)
    }


    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        
        view.findViewById<TextView>(R.id.test_text_view).text = "Robot ID: ${robotId}"
    }

}
