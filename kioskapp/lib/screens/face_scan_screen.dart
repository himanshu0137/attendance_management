import 'package:flutter/material.dart';

class FaceScanScreen extends StatefulWidget {
  const FaceScanScreen({super.key});

  @override
  State<FaceScanScreen> createState() => _FaceScanScreenState();
}

class _FaceScanScreenState extends State<FaceScanScreen> {
  bool _isScanning = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Face Scan Attendance'),
        backgroundColor: Colors.green,
        foregroundColor: Colors.white,
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(
              Icons.face_outlined,
              size: 100,
              color: Colors.green,
            ),
            const SizedBox(height: 30),
            const Text(
              'Face Recognition',
              style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 20),
            if (_isScanning)
              const CircularProgressIndicator()
            else
              const Text('Ready to scan'),
            const SizedBox(height: 30),
            ElevatedButton(
              onPressed: _isScanning ? null : _startScan,
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.green,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(horizontal: 30, vertical: 15),
              ),
              child: Text(_isScanning ? 'Scanning...' : 'Start Face Scan'),
            ),
            const SizedBox(height: 20),
            ElevatedButton(
              onPressed: () {
                // Navigate to fingerprint scan
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Fingerprint scan selected')),
                );
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.blue,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(horizontal: 30, vertical: 15),
              ),
              child: const Text('Use Fingerprint'),
            ),
          ],
        ),
      ),
    );
  }

  void _startScan() {
    setState(() {
      _isScanning = true;
    });

    // Simulate scanning process
    Future.delayed(const Duration(seconds: 3), () {
      setState(() {
        _isScanning = false;
      });
      
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Face scan completed successfully!')),
      );
    });
  }
}