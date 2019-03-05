/*
 * produce pod resource request array
 */
package main

/*
 * pod Request: cpu, mem,storage,bandwidth, typePod(hadoop,mpi,spark and others)
 * nodeName the
 */
type PodRequest struct {
	resReq   []float64
	typePod  string
	nodeName string
}
