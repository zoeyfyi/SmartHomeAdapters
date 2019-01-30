package com.github.halspals.smarthomeadapters.smarthomeadapters


import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.*
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.RobotInterface
import org.jetbrains.anko.toast

class RobotFragment : Fragment() {

    private val fTag = "RobotFragment"

    private lateinit var robotId: String
    private var robot: Robot? = null

    private lateinit var progressBar: ProgressBar
    private lateinit var switch: Switch
    private lateinit var seekBar: SeekBar

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

        progressBar = view.findViewById(R.id.progress_bar)
        switch = view.findViewById(R.id.robot_switch)
        seekBar = view.findViewById(R.id.robot_seek_bar)

        // set initial visibility
        progressBar.visibility = View.VISIBLE
        switch.visibility = View.INVISIBLE
        seekBar.visibility = View.INVISIBLE

        switch.setOnCheckedChangeListener { _, isOn ->
            onSwitch(isOn)
        }

        seekBar.setOnSeekBarChangeListener(object : SeekBar.OnSeekBarChangeListener {
            override fun onProgressChanged(seekBar: SeekBar?, progress: Int, fromUser: Boolean) {
                onSeek(progress)
            }
            override fun onStartTrackingTouch(seekBar: SeekBar?) {}
            override fun onStopTrackingTouch(seekBar: SeekBar?) {}
        })

        fetchRobot()
    }


    /**
     * Fetches the robot with id of [robotId] and calls [onReceiveRobot]
     */
    private fun fetchRobot() {
        // TODO: fetch robot from server
        // TODO: remove test code

        // odd id -> range, even id -> switch
        val isToggle = robotId.toIntOrNull()?.rem(2) == 0
        val robotInterface = if (isToggle) {
            RobotInterface.Toggle(false)
        } else {
            RobotInterface.Range(5, 0, 10)
        }

        onReceiveRobot(Robot(
            id = "123",
            nickname = "Robot 123",
            iconDrawable = R.drawable.basic_lightbulb,
            robotInterface = robotInterface
        ))
    }

    /**
     * Called whenever a new [Robot] is received
     *
     * @param robot the robot
     */
    private fun onReceiveRobot(robot: Robot) {
        this.robot = robot

        progressBar.visibility = View.INVISIBLE
        switch.visibility = View.INVISIBLE
        seekBar.visibility = View.INVISIBLE

        when(robot.robotInterface) {
            is RobotInterface.Toggle -> {
                switch.visibility = View.VISIBLE
                switch.isChecked = robot.robotInterface.isOn
            }
            is RobotInterface.Range -> {
                seekBar.visibility = View.VISIBLE
                seekBar.max = robot.robotInterface.max
                seekBar.min = robot.robotInterface.min
                seekBar.progress = robot.robotInterface.value
            }
        }
    }

    /**
     * onSwitch is called whenever the switch changes states
     *
     * @param isOn whether the switch is on/off
     */
    private fun onSwitch(isOn: Boolean) {
        Log.d(fTag, "onSwitch($isOn)")
        // TODO: send update to server
    }

    /**
     * onSeek is called whenever the seek bar changes states
     *
     * @param value value of seek bar
     */
    private fun onSeek(value: Int) {
        Log.d(fTag, "onSeek($value)")
        // TODO: send update to server
    }

}
