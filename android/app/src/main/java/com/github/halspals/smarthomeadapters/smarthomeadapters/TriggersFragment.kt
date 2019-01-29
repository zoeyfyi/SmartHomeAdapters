package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup

class TriggersFragment : Fragment() {

    private val fTag = "RobotFragment"

    override fun onCreateView(
            inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_triggers, container, false)

    }
}
