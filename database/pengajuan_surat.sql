-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Oct 19, 2023 at 09:21 AM
-- Server version: 8.0.30
-- PHP Version: 8.2.0

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `pengajuan_surat`
--

-- --------------------------------------------------------

--
-- Table structure for table `dokumen_syarat`
--

CREATE TABLE `dokumen_syarat` (
  `id` int NOT NULL,
  `id_surat` int DEFAULT NULL,
  `filename` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `dokumen_syarat`
--

INSERT INTO `dokumen_syarat` (`id`, `id_surat`, `filename`, `created_at`, `updated_at`) VALUES
(21, 19, '33090061106010001190-ktp.pdf', '2023-10-19 00:17:31', '2023-10-19 00:17:31'),
(22, 19, '33090061106010001191-ktp.pdf', '2023-10-19 00:17:31', '2023-10-19 00:17:31'),
(23, 19, '33090061106010001192-ktp.pdf', '2023-10-19 00:17:31', '2023-10-19 00:17:31'),
(24, 20, '33090061106010012200-ktp.pdf', '2023-10-19 00:17:47', '2023-10-19 00:17:47'),
(25, 20, '33090061106010012201-ktp.pdf', '2023-10-19 00:17:47', '2023-10-19 00:17:47');

-- --------------------------------------------------------

--
-- Table structure for table `masyarakat`
--

CREATE TABLE `masyarakat` (
  `idm` int NOT NULL,
  `nik` bigint DEFAULT NULL,
  `nama` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `no_hp` varchar(30) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `gender` enum('laki-laki','perempuan') COLLATE utf8mb4_general_ci DEFAULT NULL,
  `tempat_lahir` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `alamat` text COLLATE utf8mb4_general_ci,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `birthday` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `masyarakat`
--

INSERT INTO `masyarakat` (`idm`, `nik`, `nama`, `no_hp`, `gender`, `tempat_lahir`, `alamat`, `created_at`, `update_at`, `birthday`) VALUES
(1, 33090061106010001, 'Jidane Adi Ramadhzan', '082134147290', 'laki-laki', 'Boyolali', 'Tegalarum Rt09 Rw01, Sumbung, Cepogo, Boyolali', '2023-10-18 22:15:47', '2023-10-18 22:15:47', '2001-06-11'),
(2, 33090061106010002, 'Mierta Ivani Choirunnisa', '08937483913323', 'perempuan', 'Boyolali', 'Cepogo', '2023-10-18 22:17:17', '2023-10-18 22:17:17', '1999-10-22'),
(3, 33090061106010012, 'Tirta Aura Ramazan', '082134149999', 'laki-laki', 'Boyolali', 'Tegalarum Rt09 Rw01, Sumbung, Cepogo, Boyolali', '2023-10-18 22:18:21', '2023-10-18 22:18:21', '2003-12-24');

-- --------------------------------------------------------

--
-- Table structure for table `surat`
--

CREATE TABLE `surat` (
  `id` int NOT NULL,
  `id_masyarakat` int DEFAULT NULL,
  `jns_surat` enum('ktp','kematian','kelahiran','tidak mampu') COLLATE utf8mb4_general_ci DEFAULT NULL,
  `status` enum('diproses','ditolak','terverifikasi','diterbitkan','diambil') COLLATE utf8mb4_general_ci DEFAULT NULL,
  `keterangan` text COLLATE utf8mb4_general_ci,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `surat`
--

INSERT INTO `surat` (`id`, `id_masyarakat`, `jns_surat`, `status`, `keterangan`, `created_at`, `updated_at`) VALUES
(19, 1, 'ktp', 'diproses', 'Percobaan upload file', '2023-10-19 00:17:31', '2023-10-19 00:17:31'),
(20, 3, 'ktp', 'diproses', 'Percobaan upload file', '2023-10-19 00:17:47', '2023-10-19 00:17:47');

-- --------------------------------------------------------

--
-- Table structure for table `user`
--

CREATE TABLE `user` (
  `id` bigint NOT NULL,
  `email` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `role` enum('admin','masyarakat') COLLATE utf8mb4_general_ci DEFAULT NULL,
  `password` text COLLATE utf8mb4_general_ci,
  `konf_pass` text COLLATE utf8mb4_general_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `user`
--

INSERT INTO `user` (`id`, `email`, `role`, `password`, `konf_pass`) VALUES
(33090061106010001, 'jidane@gmail.com', 'masyarakat', '8618ea1de89c78504a61b41491309f60e4f4e117b209255e26d835608bc96bcd', '8618ea1de89c78504a61b41491309f60e4f4e117b209255e26d835608bc96bcd'),
(33090061106010002, 'mierta@gmail.com', 'masyarakat', '8618ea1de89c78504a61b41491309f60e4f4e117b209255e26d835608bc96bcd', '8618ea1de89c78504a61b41491309f60e4f4e117b209255e26d835608bc96bcd'),
(33090061106010012, 'tirta@gmail.com', 'masyarakat', '8618ea1de89c78504a61b41491309f60e4f4e117b209255e26d835608bc96bcd', '8618ea1de89c78504a61b41491309f60e4f4e117b209255e26d835608bc96bcd');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `dokumen_syarat`
--
ALTER TABLE `dokumen_syarat`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ID_SURAT` (`id_surat`);

--
-- Indexes for table `masyarakat`
--
ALTER TABLE `masyarakat`
  ADD PRIMARY KEY (`idm`),
  ADD KEY `NIK` (`nik`);

--
-- Indexes for table `surat`
--
ALTER TABLE `surat`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ID_MSY` (`id_masyarakat`);

--
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `dokumen_syarat`
--
ALTER TABLE `dokumen_syarat`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=26;

--
-- AUTO_INCREMENT for table `masyarakat`
--
ALTER TABLE `masyarakat`
  MODIFY `idm` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `surat`
--
ALTER TABLE `surat`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=21;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `dokumen_syarat`
--
ALTER TABLE `dokumen_syarat`
  ADD CONSTRAINT `ID_SURAT` FOREIGN KEY (`id_surat`) REFERENCES `surat` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `masyarakat`
--
ALTER TABLE `masyarakat`
  ADD CONSTRAINT `NIK` FOREIGN KEY (`nik`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `surat`
--
ALTER TABLE `surat`
  ADD CONSTRAINT `ID_MSY` FOREIGN KEY (`id_masyarakat`) REFERENCES `masyarakat` (`idm`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
