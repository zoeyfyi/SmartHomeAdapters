package com.github.halspals.smarthomeadapters.smarthomeadapters


import android.Manifest.permission.*
import android.content.BroadcastReceiver
import android.content.Context
import android.content.Context.WIFI_SERVICE
import android.content.Intent
import android.content.IntentFilter
import android.content.pm.PackageManager
import android.content.pm.PackageManager.PERMISSION_GRANTED
import android.net.wifi.ScanResult
import android.net.wifi.WifiConfiguration
import android.net.wifi.WifiManager
import android.os.Bundle
import android.support.v4.app.Fragment
import android.support.v4.content.ContextCompat.checkSelfPermission
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import com.airbnb.lottie.LottieAnimationView
import org.jetbrains.anko.alert
import org.jetbrains.anko.cancelButton
import org.jetbrains.anko.okButton
import java.util.*
import android.R.attr.start
import android.animation.ValueAnimator
import android.content.Context.CONNECTIVITY_SERVICE
import android.net.ConnectivityManager
import android.net.Network
import android.net.NetworkRequest
import android.os.Build
import android.os.Handler
import android.os.Looper
import android.view.animation.Animation
import com.airbnb.lottie.LottieDrawable


class ConnectToAdapterFragment : Fragment() {
    private val fTag = "ConnectToAdapterFrag"

    private lateinit var continueButton: Button
    private lateinit var wifiScanAnimation: LottieAnimationView

    var wifiManager: WifiManager? = null
    var scanTask: TimerTask? = null
    var scanTimer: Timer? = null

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        return inflater.inflate(R.layout.fragment_connect_to_adapter, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        continueButton = view.findViewById(R.id.continue_button)
        continueButton.isEnabled = false

        wifiScanAnimation = view.findViewById(R.id.wifi_scan_animation)
        wifiScanAnimation.setMaxFrame(WIFI_SCAN_LOADING_MAX_FRAME) // only play loading portion
        wifiScanAnimation.speed = 0.6f

        // check for wifi permission
        if (checkSelfPermission(context!!, ACCESS_WIFI_STATE) != PERMISSION_GRANTED ||
            checkSelfPermission(context!!, CHANGE_WIFI_STATE) != PERMISSION_GRANTED ||
            checkSelfPermission(context!!, ACCESS_COARSE_LOCATION) != PERMISSION_GRANTED
        ) {
            requestWifiPermissions()
        } else {
            wifiManager = setupWifiManager()
            startScanning()
        }
    }

    override fun onRequestPermissionsResult(requestCode: Int, permissions: Array<out String>, grantResults: IntArray) {
        when (requestCode) {
            CONNECT_TO_ADAPTER_REQUEST_ACCESS_WIFI_STATE -> {
                Log.d(fTag, "Got result for wifi permissions request")
                if ((grantResults.isNotEmpty() && grantResults[0] == PackageManager.PERMISSION_GRANTED)) {
                    wifiManager = setupWifiManager()
                    startScanning()
                } else {
                    // can't continue without wifi permission, quit device registration
                    activity?.finish()
                }
            }
            else -> {
            }
        }
    }

    /**
     * Starts scanning for smart home adapters
     */
    private fun startScanning() {
        val wifiManager = setupWifiManager()
        this.wifiManager = wifiManager

        // check wifi is enabled
        if (!wifiManager.isWifiEnabled) {
            Log.d(fTag, "Wifi is not enabled, enabling it")
            wifiManager.isWifiEnabled = true
        }

        // start scanning
        scanTask = object : TimerTask() {
            override fun run() {
                Log.d(fTag, "Scanning for Wifi networks")
                wifiManager.startScan()
            }
        }

        // wifi scanning is rate limited to 4 per 2 minuets, see:
        // https://developer.android.com/guide/topics/connectivity/wifi-scan#wifi-scan-throttling
        // scan every 30s to avoid throttling
        if (scanTimer == null) scanTimer = Timer()
        scanTimer?.schedule(scanTask, 0, 30 * 1000)
    }

    /**
     * Requests permissions to read and modify wifi networks
     */
    private fun requestWifiPermissions() {
        Log.d(fTag, "Requesting wifi permissions")

        if (shouldShowRequestPermissionRationale(ACCESS_WIFI_STATE)) {
            context!!.alert(
                "Wifi permissions",
                "We need to scan for Wifi devices to find your smart home adapter"
            ) {
                okButton {
                    requestWifiPermissions()
                }
                cancelButton {
                    // can't continue without wifi permission, quit device registration
                    activity?.finish()
                }
            }.show()
        } else {
            requestPermissions(
                arrayOf(ACCESS_WIFI_STATE, CHANGE_WIFI_STATE, ACCESS_COARSE_LOCATION),
                CONNECT_TO_ADAPTER_REQUEST_ACCESS_WIFI_STATE
            )
        }
    }

    /**
     * Sets up the wifi manager, with listeners for connectivity and scan results
     */
    private fun setupWifiManager(): WifiManager {
        val wifiManager = context!!.getSystemService(WIFI_SERVICE) as WifiManager
        val connectivityManager = context!!.getSystemService(CONNECTIVITY_SERVICE) as ConnectivityManager

        val networkRequest = NetworkRequest.Builder()
            .addTransportType(android.net.NetworkCapabilities.TRANSPORT_WIFI)
            .build()

        connectivityManager.registerNetworkCallback(networkRequest, object : ConnectivityManager.NetworkCallback() {
            override fun onAvailable(network: Network?) {
                Handler(Looper.getMainLooper()).post {
                    onConnectivityChange(true)
                }
            }

            override fun onUnavailable() {
                Handler(Looper.getMainLooper()).post {
                    onConnectivityChange(false)
                }
            }
        })

        // call onScanResults when we have wifi scan results
        val scanReceiver = object : BroadcastReceiver() {
            override fun onReceive(context: Context?, intent: Intent?) {
                if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
                    val success = intent?.getBooleanExtra(WifiManager.EXTRA_RESULTS_UPDATED, false)
                    Log.d(fTag, "Scan successful: $success")
                }
                onScanResults(wifiManager.scanResults)
            }
        }
        context!!.registerReceiver(scanReceiver, IntentFilter(WifiManager.SCAN_RESULTS_AVAILABLE_ACTION))

        return wifiManager
    }

    /**
     * Called whenever the wifi connectivity changes. Enables and disables the continue button.
     *
     * @param connected true if the device is connected to wifi, false otherwise
     */
    private fun onConnectivityChange(connected: Boolean) {
        Log.d(fTag, "Connectivity changed, connected: $connected")

        // check we are connected to the adapter
        val connectedToAdapter = connected && wifiManager?.let { isNetworkAnAdapter(it.connectionInfo.ssid) } ?: false
        Log.d(fTag, "Connected to adapter: $connectedToAdapter")

        if (connectedToAdapter) {
            continueButton.isEnabled = true
            wifiScanAnimation.setMaxFrame(Int.MAX_VALUE)
            wifiScanAnimation.loop(false)
        } else {
            continueButton.isEnabled = false
            wifiScanAnimation.setMaxFrame(WIFI_SCAN_LOADING_MAX_FRAME)
            wifiScanAnimation.loop(true)
            wifiScanAnimation.playAnimation() // restart animation
        }
    }

    /**
     * Called whenever a network scan is finished
     *
     * @param results list of networks from wifi scan
     */
    private fun onScanResults(results: List<ScanResult>) {
        Log.d(fTag, "Found ${results.size} Wifi networks")

        // filter smart home adapters
        val adapters = results.filter { isNetworkAnAdapter(it.SSID) }
        Log.d(fTag, "Found ${adapters.size} adapters")

        // try to the first adapter
        if (adapters.isNotEmpty()) {
            scanTimer?.cancel() // stop scanning
            connectToAdapter(adapters[0].SSID)
        }
    }

    /**
     * isNetworkAnAdapter checks if the SSID of a network matches the format of a smart home adapter
     *
     * @param ssid SSID of the network
     * @return true if the SSID matches an adapter
     */
    private fun isNetworkAnAdapter(ssid: String): Boolean {
        return ssid.contains("Smart Home Adapter")
    }

    /**
     * connectToAdapter connects to the smart home adapters temporary wifi network
     *
     * @param ssid SSID of the adapters network
     */
    private fun connectToAdapter(ssid: String) {
        Log.d(fTag, "Connecting to adapter with SSID \"$ssid\"")

        // add wifi network
        val configuration = WifiConfiguration().apply {
            SSID = ssid
            allowedKeyManagement.set(WifiConfiguration.KeyMgmt.NONE)
        }

        // connect to wifi network
        wifiManager?.apply {
            val networkId = addNetwork(configuration)
            disconnect()
            enableNetwork(networkId, true)
            reconnect()
        }
    }

    companion object {
        private const val CONNECT_TO_ADAPTER_REQUEST_ACCESS_WIFI_STATE = 1000;
        private const val WIFI_SCAN_LOADING_MAX_FRAME = 50
    }

}
