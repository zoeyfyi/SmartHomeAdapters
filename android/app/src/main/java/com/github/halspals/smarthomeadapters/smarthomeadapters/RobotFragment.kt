package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.*
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import com.sdsmdg.harjot.crollerTest.Croller
import com.sdsmdg.harjot.crollerTest.OnCrollerChangeListener
import it.beppi.tristatetogglebutton_library.TriStateToggleButton
import kotlinx.android.synthetic.main.fragment_robot.*
import okhttp3.ResponseBody
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.toast
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class RobotFragment : Fragment() {

    private val fTag = "RobotFragment"

    private lateinit var parent: MainActivity

    private lateinit var robotId: String
    private var robot: Robot? = null

    private lateinit var progressBar: ProgressBar
    private lateinit var switch: TriStateToggleButton
    private lateinit var seekBar: Croller

    private var intermediateColor: Int? = null
    private var finishedColor: Int? = null

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {

        // get robotId from bundle
        val robotIdArgument = arguments?.getString("robotId")
        if (robotIdArgument == null) {
            // no id passed, try to go back
            Log.d(fTag, "No robotId passed to robotFragment")
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

        parent = activity as MainActivity

        progressBar = view.findViewById(R.id.progress_bar)
        switch = view.findViewById(R.id.robot_switch)
        seekBar = view.findViewById(R.id.robot_seek_bar)

        intermediateColor = view.context.getColor(R.color.colorIntermediate)
        finishedColor = view.context.getColor(R.color.colorOn)

        // set initial visibility
        progressBar.visibility = View.VISIBLE
        switch.visibility = View.INVISIBLE
        seekBar.visibility = View.INVISIBLE

        fetchRobot()
    }


    /**
     * Fetches the robot with id of [robotId] and calls [onReceiveRobot]
     */
    private fun fetchRobot() {
        parent.restApiService
                .getRobot(robotId, parent.authToken)
                .enqueue(object: Callback<Robot> {

            override fun onResponse(call: Call<Robot>, response: Response<Robot>) {
                val robot = response.body()

                if (response.isSuccessful && robot != null) {
                    Log.d(fTag, "Successfully retrieved $robot")
                    onReceiveRobot(robot)
                } else {
                    val error = RestApiService.extractErrorFromResponse(response)
                    Log.e(fTag, "Getting the robot was unsuccessful, error: $error")
                    if (error != null) {
                        snackbar_layout.snackbar(error)
                    }
                }
            }

            override fun onFailure(call: Call<Robot>, t: Throwable) {
                val error = t.message
                Log.e(fTag, "fetchRobot() failed: $error")
                if (error != null) {
                    snackbar_layout.snackbar(error)
                }
            }
        })
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

        when(robot.robotInterfaceType) {
            Robot.TYPE_TOGGLE -> {
                switch.visibility = View.VISIBLE
                switch.toggleStatus = if (robot.robotStatus.value) {
                    TriStateToggleButton.ToggleStatus.on
                } else {
                    TriStateToggleButton.ToggleStatus.off
                }

                switch.setOnToggleChanged { toggleStatus, _, _ ->
                    Log.d(fTag, "[switch onToggle]: $toggleStatus")

                    when (toggleStatus) {
                        null -> Log.e(fTag, "Switch got null toggleStatus")
                        TriStateToggleButton.ToggleStatus.on -> {
                            onSwitch(
                                    true,
                                    TriStateToggleButton.ToggleStatus.on,
                                    TriStateToggleButton.ToggleStatus.off)
                            switch.toggleStatus = TriStateToggleButton.ToggleStatus.mid
                        }
                        TriStateToggleButton.ToggleStatus.off -> {
                            onSwitch(
                                    false,
                                    TriStateToggleButton.ToggleStatus.off,
                                    TriStateToggleButton.ToggleStatus.on)
                            switch.toggleStatus = TriStateToggleButton.ToggleStatus.mid
                        }
                        TriStateToggleButton.ToggleStatus.mid -> {
                            Log.e(fTag, "Switch should not get event for toggle set to mid")
                        }
                    }
                }
            }

            Robot.TYPE_RANGE -> {
                seekBar.max = robot.robotStatus.max - robot.robotStatus.min
                seekBar.progress = robot.robotStatus.current - robot.robotStatus.min
                seekBar.label = robot.robotStatus.current.toString()
                seekBar.visibility = View.VISIBLE

                seekBar.setOnCrollerChangeListener(object : OnCrollerChangeListener {
                    override fun onProgressChanged(croller: Croller?, progress: Int) {}
                    override fun onStartTrackingTouch(croller: Croller?) {}
                    override fun onStopTrackingTouch(croller: Croller?) {
                        if (croller == null) {
                            Log.w(fTag, "[onStopTrackingTouch] got null Croller")
                        } else {
                            val seekValue = croller.progress + robot.robotStatus.min
                            onSeek(seekValue)
                        }
                    }
                })
            }

            else -> {
                TODO("No other robot interface types exepcted")
            }
        }
    }

    /**
     * onSwitch is called whenever the switch changes states
     *
     * @param isOn whether the switch is on/off
     * @param goToSuccess the toggle status to go to for the switch on success
     * @param goToFailure the toggle status to go to for the switch on failure
     */
    private fun onSwitch(
            isOn: Boolean,
            goToSuccess: TriStateToggleButton.ToggleStatus,
            goToFailure: TriStateToggleButton.ToggleStatus)
    {
        Log.d(fTag, "onSwitch($isOn)")
        progressBar.visibility = View.VISIBLE

        // Send the update to the server
        parent.restApiService
                .robotToggle(robotId, isOn, parent.authToken, mapOf())
                .enqueue(object: Callback<ResponseBody> {

            override fun onResponse(call: Call<ResponseBody>, response: Response<ResponseBody>) {
                progressBar.visibility = View.INVISIBLE

                if (response.isSuccessful) {
                    parent.toast("Success")
                    Log.d(fTag, "Server accepted setting switch to $isOn")
                    switch.toggleStatus = goToSuccess
                } else {
                    val error = RestApiService.extractErrorFromResponse(response)
                    Log.e(fTag, "Setting the switch was unsuccessful, error: $error")
                    switch.toggleStatus = goToFailure
                    if (error != null) {
                        snackbar_layout.snackbar(error)
                    }
                }
            }

            override fun onFailure(call: Call<ResponseBody>, t: Throwable) {
                progressBar.visibility = View.INVISIBLE
                val error = t.message
                Log.e(fTag, "onSwitch($isOn) FAILED, error: $error")
                switch.toggleStatus = goToFailure
                if (error != null) {
                    snackbar_layout.snackbar(error)
                }
            }
        })
    }

    /**
     * onSeek is called whenever the seek bar changes states
     *
     * @param value current of seek bar
     */
    private fun onSeek(value: Int) {
        Log.d(fTag, "onSeek($value)")
        progressBar.visibility = View.VISIBLE

        seekBar.label = value.toString()
        seekBar.isEnabled = false
        val startColor = intermediateColor
        if (startColor != null) { seekBar.labelColor = startColor }
        val endColor = finishedColor

        // TODO make more elegant solution than just going K->C by removing 273
        parent.restApiService
                .robotRange(robotId, value-273, parent.authToken, mapOf())
                .enqueue(object : Callback<ResponseBody> {

            override fun onResponse(call: Call<ResponseBody>, response: Response<ResponseBody>) {
                // Indicate that the async call has finished
                progressBar.visibility = View.INVISIBLE
                seekBar.isEnabled = true
                if (endColor != null) { seekBar.labelColor = endColor }

                // Handle the callback
                if (response.isSuccessful) {
                    parent.toast("Success")
                    Log.d(fTag, "Server accepted setting range to $value")
                } else {
                    val error = RestApiService.extractErrorFromResponse(response)
                    Log.e(fTag, "Setting the range was unsuccessful, error: $error")
                    if (error != null) {
                        snackbar_layout.snackbar(error)
                    }
                }
            }

            override fun onFailure(call: Call<ResponseBody>, t: Throwable) {
                // Indicate that the async call has finished
                progressBar.visibility = View.INVISIBLE
                seekBar.isEnabled = true
                if (endColor != null) { seekBar.labelColor = endColor }

                // Handle the error throwable
                val error = t.message
                Log.e(fTag, "onSeek($value) FAILED, error: $error")
                if (error != null) {
                    snackbar_layout.snackbar(error)
                }
            }
        })
    }

}
