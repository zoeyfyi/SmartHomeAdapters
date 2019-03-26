package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import kotlinx.android.synthetic.main.fragment_manual_entry.*

/**
 * A fragment which allows the user to manually enter the robot's registration code.
 */
class ManualEntryFragment : Fragment() {

    private val fTag = "ManualEntryFragment"

    private val parent by lazy { activity as RegisterRobotActivity }

    override fun onCreateView(
            inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?
    ): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_manual_entry, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)


        continue_button.setOnClickListener { _ ->
            parent.robotId = registration_code.text.toString()
            parent.startFragment(NicknameFragment())
        }

        cancel_button.setOnClickListener { _ -> parent.startFragment(QRFragment()) }

    }
}
