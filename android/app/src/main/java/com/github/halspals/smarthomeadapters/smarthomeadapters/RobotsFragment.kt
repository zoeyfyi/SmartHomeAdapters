package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.GridView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import kotlinx.android.synthetic.main.fragment_robots.*
import org.jetbrains.anko.design.snackbar
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class RobotsFragment : Fragment() {

    private val fTag = "RobotsFragment"

    private lateinit var robotGrid: GridView

    private val parent by lazy { activity as MainActivity }

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?
    ): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_robots, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        parent.authState.performActionWithFreshTokens(parent.authService)
        { accessToken, _, ex ->
            if (accessToken == null) {
                Log.e(fTag, "[performActionWithFreshTokens] got null access token, "
                        + "exception: $ex")
            } else {
                fetchRobots(accessToken, view)
            }
        }
    }

    private fun fetchRobots(token: String, view: View) {
        parent.restApiService
                .getRobots(token)
                .enqueue(object : Callback<List<Robot>> {

                    override fun onFailure(call: Call<List<Robot>>, t: Throwable) {
                        val errorMsg = t.message
                        Log.e(fTag, "getRobots FAILED, got error: $errorMsg")
                        if (errorMsg != null) {
                            snackbar_layout.snackbar(errorMsg)
                        }
                    }

                    override fun onResponse(
                            call: Call<List<Robot>>,
                            response: Response<List<Robot>>) {

                        val robots = response.body()
                        if (response.isSuccessful && robots != null) {
                            displayRobots(view, robots)
                        } else {
                            val error = RestApiService.extractErrorFromResponse(response)

                            Log.e(fTag, "getRobots got unsuccessful response, error: $error")
                            if (error != null) {
                                snackbar_layout.snackbar(error)
                            }
                        }
                    }
                })
    }

    private fun displayRobots(view: View, robots: List<Robot>) {
        robotGrid = view.findViewById(R.id.RobotGrid)
        robotGrid.adapter = RobotAdapter(view.context, robots.toMutableList()) { robot ->
            Log.d(fTag, "Clicked robot: \"${robot.nickname}\"")

            // create fragment with robot ID
            val robotFragment = RobotFragment()
            val bundle = Bundle()
            bundle.putString("robotId", robot.id)
            bundle.putString("robotType", robot.robotType)
            robotFragment.arguments = bundle

            parent.startFragment(robotFragment, true)
        }
    }
}
